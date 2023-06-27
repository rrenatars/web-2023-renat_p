const postName = document.getElementById('PostName');
const postDescription = document.getElementById('PostDescription');
const authorName = document.getElementById('AuthorName');
const publishDate = document.getElementById('PublishDate');
const postContent = document.getElementById('PostContent');

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

let avatarUrlBase64 = "";
let avatarFileName = "";

function previewAuthorAvatar() {
    const file = authorAvatar.files[0];

    if (file) {
        const reader = new FileReader();

        reader.addEventListener('load', function() {
            avatarBase64 = reader.result;
            avatarUrlBase64 = avatarBase64.replace("data:", "").replace(/^.+,/, "")

            const elements = document.querySelectorAll('.author__avatar, .photo__area');
            Array.prototype.forEach.call(elements, el => {
                const imgElement = document.createElement('img');
                imgElement.src = avatarBase64;
                avatarFileName = file.name
                imgElement.setAttribute('class', 'area__avatar_preview');

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

let heroimageBigUrlBase64 = "";
let heroimageSmallUrlBase64 = "";
let bigHeroimageFileName = "";
let smallHeroimageFileName = "";

function previewImage(image, elements) {
    const file = image.files[0];

    if (file) {
        const reader = new FileReader();

        reader.addEventListener('load', function() {
            const heroimageBase64 = reader.result;
            if (elements[1].id === 'PreviewBigImage') {
                heroimageBigUrlBase64 = heroimageBase64.replace("data:", "").replace(/^.+,/, "");
                bigHeroimageFileName = file.name;
            }
            if (elements[1].id === 'PreviewSmallImage') {
                heroimageSmallUrlBase64 = heroimageBase64.replace("data:", "").replace(/^.+,/, "");
                smallHeroimageFileName = file.name
            }

            Array.prototype.forEach.call(elements, el => {
                const imgElement = document.createElement('img');
                imgElement.src = heroimageBase64;
                imgElement.setAttribute('class', 'photo__preview');

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

function deleteImage(el) {
    let temp = el.closest("p");
    if (temp.id != 'UploadButton') {
        if (temp.id === 'SmallHeroimageHint') {
            const heroimageHintSmall = document.querySelector('.heroimage__hint_smaller');
            deleteHeroimage(heroimageHintSmall);
        }
        if (temp.id === 'BigHeroimageHint') {
            const heroimageHintBig = document.querySelector('.heroimage__hint');
            deleteHeroimage(heroimageHintBig);
        }
    } else {
        deleteAuthorAvatar(temp);
    }
}

function deleteAuthorAvatar(element) {
    let hintText;
    hintText = avatarHintText;
    document.getElementById(element.id).textContent = hintText;
    const avatarArea = document.querySelector('.area__avatar_preview');
    avatarArea.src = '/static/images/camera.svg';
    avatarArea.setAttribute('class', 'area__avatar');
    const previewAvatar = document.querySelector('.author__avatar');
    const previewAvatarImage = previewAvatar.firstElementChild;
    previewAvatar.removeChild(previewAvatarImage);
}

function deleteHeroimage(element) {
    let elementId = element.id;
    let hintText;
    let selector;
    let preview;
    let hint;
    if (elementId === 'BigHeroimageHint') {
        hint = bigHeroimageHint;
        hintText = bigHeroimageHintText;
        selector = '.photo__heroimage';
        preview = previewBigImage;
    }
    if (elementId === 'SmallHeroimageHint') {
        hint = smallHeroimageHint;
        hintText = smallHeroimageHintText;
        selector = '.photo__heroimage_smaller';
        preview = previewSmallImage;
    }
    document.getElementById(elementId).textContent = hintText;
    const container = document.querySelector(selector);
    const imgElement = container.firstElementChild;
    imgElement.src = '/static/images/camera.svg';
    imgElement.setAttribute('class', 'heroimage__avatar');
    const spanElement = document.createElement('span');
    spanElement.textContent = 'Upload';
    spanElement.setAttribute('class', 'buttons__button-new');
    container.appendChild(spanElement);
    const container2 = preview;
    const imgElement2 = container2.firstElementChild;
    container2.removeChild(imgElement2);
}

document.body.addEventListener('click', function (e) {
    if (e.target.className === 'label-disable') {
        e.preventDefault();
    }
});

const publishButton = document.getElementById('Submit');

publishButton.addEventListener("click", () => {
    const formData = {
        title: postName.value,
        description: postDescription.value,
        author_name: authorName.value,
        author_avatar: avatarUrlBase64,
        avatar_file_name: avatarFileName,
        publish_date: publishDate.value,
        big_heroimage: heroimageBigUrlBase64,
        big_heroimage_file_name: bigHeroimageFileName,
        small_heroimage: heroimageSmallUrlBase64,
        small_heroimage_file_name: smallHeroimageFileName,
        content: postContent.value
    }

    console.log(JSON.stringify(formData, null, "\t"));
    createPost(formData);
  }    
)

async function createPost(data) {
    const response = await fetch("/api/post", {
      method: "POST",
      headers: {
        "Content-Type": "application/json;charset=utf-8"
      },
      body: JSON.stringify(data)
    });
    if (!response.ok) {
      alert("Ошибка HTTP: " + response.status);
    }
  }







