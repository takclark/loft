var listEndpoint = "/list";
var slides, slideInterval, imageUpdateInterval;
var showingSlide = 0;
var imageList;
var defaultImages = ['1521062876586560000.jpeg', '1521062790904012000.jpeg'];
var imageQueue = {
  'imageList': defaultImages
}
var imgTagTemplate = '<img src="/images/%FNAME%">';

function getNextImage() {
  let nextFilename = imageQueue.imageList.shift();
  imageQueue.imageList[imageQueue.imageList.length] = nextFilename;
  return nextFilename;
}

function getUpdatedImageList() {
  $.getJSON(listEndpoint, updateImageList);
}

function updateImageList(listData, textStatus) {
  if (textStatus != 'success') {
    console.log('bad status on list response, return');
    return;
  }

  let images = listData.images
  // first build list of image names
  let imageFilesOnServer = []
  for (let i = 0; i < images.length; i++) {
    imageFilesOnServer.push(images[i].filename);
  }

  // remove anything in the local list that's not on the server
  imageQueue.imageList = imageQueue.imageList.filter(name => imageFilesOnServer.includes(name));

  // take everything out of the sever's list that's already in ours
  imageFilesOnServer = imageFilesOnServer.filter(name => !imageQueue.imageList.includes(name));

  // anything we didn't have would be new, so push it to the front of our local image list
  imageQueue.imageList = imageFilesOnServer.concat(imageQueue.imageList);
}

function getImageTagFromFilename(fname) {
  return imgTagTemplate.replace('%FNAME%', fname);
}

function init() {
  slides = document.querySelectorAll('#slides .slide');
  slideInterval = setInterval(nextSlide, 5000);
  imageUpdateInterval = setInterval(getUpdatedImageList, 10000);
}

// update the image on the nextSlide and show it
function nextSlide() {
  prevSlide = showingSlide
  slides[showingSlide].className = 'slide';
  showingSlide = (showingSlide+1)%2;
  slides[showingSlide].className = 'slide showing';

  // update the image src on the now-hidden slide
  setTimeout(function() {updateSlideElementImage(slides[prevSlide])}, 1000);
}

function updateSlideElementImage(elem) {
  elem.innerHTML = getImageTagFromFilename(getNextImage());
}

$( document ).ready(init)
