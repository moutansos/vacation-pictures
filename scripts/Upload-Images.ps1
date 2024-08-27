param(
    # Vacation ID
    [Parameter(Mandatory)]
    [string]
    $VacationId,

    # The path to the folder to upload
    [Parameter(Mandatory)]
    [string]
    $PathToUpload,

    # The tags to attach to each file
    [Parameter(Mandatory=$false)]
    [string[]]
    $Tags = @()
)

$tempDirectory = "/tmp"

function Generate-Thumbnail($id, $imagePath) {
    Write-Host "Starting image generation for $id`: $imagePath"
    $newFullFileName = "$id-thumb.webp"
    $tempFileLocation = "$tempDirectory/$newFullFileName"
    $imageFile = Get-ChildItem $imagePath
    $extension = $imageFile.Extension

    $resizedTempFile = "$tempDirectory/$id-resized.$extension"
    Write-Host "Resizing..."
    convert $imagePath -thumbnail 3000x250 $resizedTempFile

    Write-Host "Converting to webp..."
    cwebp -q 80 -o $tempFileLocation $resizedTempFile

    WRItE-Host "Uploading to s3..."
    aws s3 cp $tempFileLocation "s3://msyke-vaca-pictures/$newFullFileName"

    Write-Host "Finished with $id`: $imagePath"
    return "https://msyke-vaca-pictures.s3.amazonaws.com/$newFullFileName"
}

function Generate-Image([guid]$id, $imagePath) {
    Write-Host "Starting image generation for $id`: $imagePath"
    $newFullFileName = "$id.webp"
    $tempFileLocation = "$tempDirectory/$newFullFileName"

    Write-Host "Converting to webp..."
    cwebp -q 80 -o $tempFileLocation $imagePath 

    Write-Host "Uploading to s3..."
    aws s3 cp $tempFileLocation "s3://msyke-vaca-pictures/$newFullFileName"

    Write-Host "Finished with $id`: $imagePath"
    return "https://msyke-vaca-pictures.s3.amazonaws.com/$newFullFileName"
}

$imageFiles = Get-ChildItem $PathToUpload -File
$imageFiles = $imageFiles | Where-Object { $_.Extension -match "jpg|jpeg|png" }

$vacationsFile = Resolve-Path "$PSScriptRoot/../vacations.json"
[string]$vacationsDataString = Get-Content $vacationsFile -Raw
[pscustomobject]$vacationsData = $vacationsDataString | ConvertFrom-Json

[pscustomobject[]]$vacations = $vacationsData.vacations
[pscustomobject[]]$selectedVacations = $vacations | Where-Object { $_.id -eq $VacationId }

if($selectedVacations.Length -eq 0) {
    Write-Host "Didn't find any matching vacations"
    return
}

$selectedVacation = $selectedVacations[0]


foreach ($img in $imageFiles) {
    $path = $img.FullName
    $fileName = $img.Name
    $id = [guid]::NewGuid()
    $newPic = [PSCustomObject]@{
        title = $fileName;
        description = "";
        thumbnailPath = (Generate-Thumbnail $id $path)[-1];
        imagePath = (Generate-Image $id $path)[-1];
        tags = $Tags;
    }
    $selectedVacation.pictures = $selectedVacation.pictures + @($newPic)
        $newVacationDataString = $vacationsData | ConvertTo-Json -Depth 30
        $newVacationDataString | Set-Content -Path $vacationsFile
}

