function scrollToSelectedThumbnail() {
    const currentSelection = document.querySelector('.selected-thumbnail');
    if (currentSelection) {
        // if (currentSelection.scrollIntoViewIfNeeded) {
        //     currentSelection.scrollIntoViewIfNeeded();
        //     return;
        // }
        currentSelection.scrollIntoView({ behavior: "smooth" });
    }
}
document.addEventListener('htmx:afterSwap', function (event) {
    if (event.detail.target.id !== 'view-box') {
        return;
    }

    /** @type {string?} */
    const previousIndex = event.detail.target.getAttribute('data-idx');
    const previousSelectionId = `#thumbnail-${previousIndex}`;    

    const prevImage = document.querySelector(previousSelectionId);
    if (prevImage) {
        prevImage.classList.remove('selected-thumbnail');
    }

    setTimeout(() => {
        scrollToSelectedThumbnail();
    }, 100);
});

document.addEventListener('DOMContentLoaded', function () {
    const fullScreenButton = document.getElementById('full-screen');
    if (fullScreenButton) {
        fullScreenButton.addEventListener('click', enterFullscreen);
    }
});



let isFullScreen = false;
async function enterFullscreen() {
    const picViewer = document.querySelector('.pic-viewer');
    if (!picViewer) {
        return;
    }

    if (!isFullScreen) {
        if (picViewer?.requestFullscreen) {
            await picViewer.requestFullscreen();
            picViewer.classList.add('full-screen-pic-viewer');
            isFullScreen = true;
        }
    } else {
        if (document.exitFullscreen) {
            document.exitFullscreen();
            picViewer.classList.remove('full-screen-pic-viewer');
            isFullScreen = false;
        }
    }
}

window.enterFullScreenViewBox = enterFullscreen;
