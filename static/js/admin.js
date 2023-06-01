const postName = document.getElementById('PostName');
const postDescription = document.getElementById('PostDescription');
const authorName = document.getElementById('AuthorName');
const publishDate = document.getElementById('PublishDate');

postName.addEventListener('keyup', function () {
    update(postName, 'content__title');
});

postDescription.addEventListener('keyup', function () {
    update(postDescription, 'content__description');
});

authorName.addEventListener('keyup', function () {
    update(authorName, 'author__name');
});

publishDate.addEventListener('change', function () {
    update(publishDate, 'info__date');
});

function update(element, previewClass) {
    const input = element.value;
    const classes = document.getElementsByClassName(previewClass);
    Array.prototype.forEach.call(classes, el => {
        el.innerText = input;
    });
}

const authorAvatar = document.getElementById('AuthorAvatar');
const uploadButton = document.getElementById('UploadButton');
const buttons = document.getElementById('Buttons');

authorAvatar.addEventListener('change', function () {
    previewAuthorAvatar();
    replace(uploadButton, buttons);
});

function previewAuthorAvatar() {
    const file = authorAvatar.files[0];

    if (file) {
        const reader = new FileReader();

        reader.addEventListener('load', function() {
            const imageUrl = reader.result;

            const classes = document.querySelectorAll('.author__avatar, .photo__area');
            Array.prototype.forEach.call(classes, el => {
                const imgElement = document.createElement('img');
                imgElement.src = imageUrl;
                imgElement.style.width = '100%';
                imgElement.style.height = '100%';
                imgElement.style.borderRadius = '50%';

                el.innerHTML = '';
                el.appendChild(imgElement);
            });
        });

        reader.readAsDataURL(file);
    }
}

const heroimageBig = document.getElementById('HeroimageBig');
const classesForBigImage = document.querySelectorAll('.photo__heroimage, .content__picture');
const bigHeroimageHint = document.getElementById('BigHeroimageHint');

heroimageBig.addEventListener('change', function () {
    previewImage(heroimageBig, classesForBigImage);
    replace(bigHeroimageHint, buttons);
});

const heroimageSmall = document.getElementById('HeroimageSmall');
const classesForSmallImage = document.querySelectorAll('.photo__heroimage_smaller, .content__picture_home');
const smallHeroimageHint = document.getElementById('SmallHeroimageHint');

heroimageSmall.addEventListener('change', function () {
    previewImage(heroimageSmall, classesForSmallImage);
    replace(smallHeroimageHint, buttons);
});

function previewImage(image, classes) {
    const file = image.files[0];

    if (file) {
        const reader = new FileReader();

        reader.addEventListener('load', function() {
            const imageUrl = reader.result;

            Array.prototype.forEach.call(classes, el => {
                const imgElement = document.createElement('img');
                imgElement.src = imageUrl;
                imgElement.setAttribute('class', 'photo__preview')

                el.innerHTML = '';
                el.appendChild(imgElement);
            });
        });

        reader.readAsDataURL(file);
    }
}

function replace(hint, buttons) {
    hint.innerHTML = buttons.innerHTML;
}

const bigHeroimageHintText = bigHeroimageHint.innerText;
const previewBigImage = document.getElementById('PreviewBigImage');

const smallHeroimageHintText = smallHeroimageHint.innerText;
const previewSmallImage = document.getElementById('PreviewSmallImage');

const avatarHint = document.getElementById('UploadButton');
const avatarHintText = avatarHint.innerText;
const avatarPreview = document.getElementsByClassName('.photo__area');

function deleteImage(el) {
    let temp = el.closest("p"); ///query поменять
    let hintText;
    let selector;
    let preview;
    let hint;
    if (temp.id === 'UploadButton') {
        hint = avatarHint;
        hintText = avatarHintText;
        selector = document.querySelectorAll('.photo__area, .author__avatar');
        document.getElementById(temp.id).textContent = hintText;
        Array.prototype.forEach.call(selector, container => {
            const imgElement = container.firstElementChild;
            imgElement.src = '../static/images/camera.svg';
            imgElement.style.width = '24px';
            imgElement.style.height = '24px';
            imgElement.style.borderRadius = '70%';
        });
    } else {
        if (temp.id === 'BigHeroimageHint') {
            hint = bigHeroimageHint;
            hintText = bigHeroimageHintText;
            selector = '.photo__heroimage';
            preview = previewBigImage;
        }
        if (temp.id === 'SmallHeroimageHint') {
            hint = smallHeroimageHint;
            hintText = smallHeroimageHintText;
            selector = '.photo__heroimage_smaller';
            preview = previewSmallImage;
        }
        document.getElementById(temp.id).textContent = hintText;
        const container = (temp.closest("div")).querySelector(selector);
        const imgElement = container.firstElementChild;
        imgElement.src = '../static/images/camera.svg'; // из корня
        imgElement.setAttribute('class', 'heroimage__avatar');
        const spanElement = document.createElement('span');
        spanElement.textContent = 'Upload';
        spanElement.setAttribute('class', 'buttons__button-new');
        container.appendChild(spanElement);
        const container2 = (preview.closest("div"))
        const imgElement2 = container2.firstElementChild;
        container2.removeChild(imgElement2);
    }

}

document.body.addEventListener('click', function (e) {
    if (e.target.className === 'label-disable') {
        e.preventDefault();
    }
});

const publishButton = document.querySelector('.publish__button');

publishButton.addEventListener('click', function() {
    const formData = new FormData();

    const inputElements = document.querySelectorAll('input, textarea');

    inputElements.forEach((input) => {
        const name = input.getAttribute('name');
        const value = input.value;

        formData.append(name, value);
    });


    const jsonData = JSON.stringify(Object.fromEntries(formData));

    console.log(jsonData);
});




