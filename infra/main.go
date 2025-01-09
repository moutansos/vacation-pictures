package main

import (
	"fmt"
	"vacation-pictures-infra/bw"

	"github.com/moutansos/vacation-pictures/common"

	kubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const APP_NAME = "vacationpictures"
const INTERNAL_CONTAINER_PORT = 8081

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackName := ctx.Stack()
		containerTag, ctExists := ctx.GetConfig("vacation-pictures-infra:containerTag")
		fmt.Printf("ContainerTag: %s\n", containerTag)
		if !ctExists {
			return fmt.Errorf("The pulumi config containerTag does not exist!")
		}

		containerImageName := fmt.Sprintf("ghcr.io/moutansos/vacationpictures:%s", containerTag)

		bwCtx, err := bw.InitBw()
		if err != nil {
			return err
		}
		defer bwCtx.Close()

		slackWebhookUrl, err := bwCtx.BuildEnvVar(stackName, common.VACA_SLACK_WEBHOOK_URL)
		if err != nil {
			return err
		}

		kubeConfig, err := bwCtx.GetByEnvAndName(stackName, common.KUBECONFIG)
		if err != nil {
			return err
		}

		fmt.Println("Creating Kubernetes Provider...")
		k8sProvider, err := kubernetes.NewProvider(ctx, "msykek8s1", &kubernetes.ProviderArgs{
			Kubeconfig: pulumi.String(kubeConfig),
		})
		if err != nil {
			return err
		}

		_, err = appsv1.NewDeployment(ctx, APP_NAME, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(APP_NAME),
				Labels: pulumi.StringMap{
					"app": pulumi.String(APP_NAME),
				},
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(2),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String(APP_NAME),
					},
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String(APP_NAME),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: &corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String(containerImageName),
								Name:  pulumi.String(APP_NAME),
								Ports: &corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(INTERNAL_CONTAINER_PORT),
									},
								},
								Env: &corev1.EnvVarArray{
									slackWebhookUrl,
									&corev1.EnvVarArgs{
										Name:  pulumi.String("APP_ENV"),
										Value: pulumi.String(stackName),
									},
								},
							},
						},
						ImagePullSecrets: &corev1.LocalObjectReferenceArray{
							&corev1.LocalObjectReferenceArgs{
								Name: pulumi.String("ghcr-secret"),
							},
						},
					},
				},
			},
		}, pulumi.Provider(k8sProvider))

		if err != nil {
			return err
		}

		_, err = corev1.NewService(ctx, APP_NAME, &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(APP_NAME),
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: &corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
						TargetPort: pulumi.Int(INTERNAL_CONTAINER_PORT),
					},
				},
				Selector: pulumi.StringMap{
					"app": pulumi.String(APP_NAME),
				},
				Type: pulumi.String("ClusterIP"),
			},
		}, pulumi.Provider(k8sProvider))

		if err != nil {
			return err
		}

		return nil
	})
}
