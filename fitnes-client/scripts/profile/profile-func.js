export function addHeader(headerText){
    const userInfoElement = document.getElementById('user-info');
    const h1Element = document.createElement('h1');
    h1Element.textContent = headerText;
    userInfoElement.insertBefore(h1Element, userInfoElement.firstChild);
}
export function loadUserProfile(token, tokenPayload) {
    try {
        // Извлекаем userid из токена
        const userId = tokenPayload.uid;

        // Выполняем GET запрос к /account/profile с указанием userid
        fetch('http://localhost:8080/account/profile', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ userid: userId })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при выполнении запроса: ' + response.statusText);
                }
                return response.json();
            })
            .then(data => {
                // Обработка данных профиля
                console.log(data);
                setUserData(data)

                // Возможно, здесь будет код для отображения данных профиля на странице
            })
            .catch(error => {
                console.error('Ошибка при выполнении запроса:', error);
            });
    } catch (error) {
        console.error('Ошибка при декодировании токена:', error);
    }
}
function setUserData(userData){
    // Получаем элементы, куда будем вставлять информацию о пользователе
    var nameElement = document.getElementById("user-name");
    var emailElement = document.getElementById("user-login");
    var idElement = document.getElementById("user-id");
    var patronymicElement = document.getElementById("user-patronymic");
    var phoneNumberElement = document.getElementById("user-phone");
    var roleElement = document.getElementById("user-role");
    var surnameElement = document.getElementById("user-surname");


    // Вставляем данные о пользователе в соответствующие элементы
    nameElement.textContent = userData.name;
    emailElement.textContent = userData.email;
    idElement.textContent = userData.id;
    patronymicElement.textContent = userData.patronymic;
    phoneNumberElement.textContent = userData.phoneNumber;
    roleElement.textContent = userData.role;
    surnameElement.textContent = userData.surname;
}

export function createButtonWithTextAndHandler(buttonText, handler)  {
    const htmlButtonElement = document.createElement("button");
    htmlButtonElement.textContent = buttonText;
    htmlButtonElement.classList.add("btn", "btn-primary", "ml-2"); // Добавляем классы Bootstrap

    // Добавление обработчика события для кнопки
    htmlButtonElement.addEventListener('click',handler)
    return htmlButtonElement
}
export function parseUserData() {
    // Создаем объект для хранения данных пользователя
    const userData = {};

    // Получаем элементы DOM для каждого поля
    const userName = document.getElementById('user-name');
    const userSurname = document.getElementById('user-surname');
    const userPatronymic = document.getElementById('user-patronymic');
    const userLogin = document.getElementById('user-login');
    const userId = document.getElementById('user-id');
    const userPhone = document.getElementById('user-phone');
    const userRole = document.getElementById('user-role');

    // Заполняем объект данными
    userData.name = userName.textContent;
    userData.surname = userSurname.textContent;
    userData.patronymic = userPatronymic.textContent;
    userData.login = userLogin.textContent;
    userData.id = userId.textContent;
    userData.phone = userPhone.textContent;
    userData.role = userRole.textContent;

    // Возвращаем объект с данными
    return userData;
}




export function addTitle(titleText) {
    const infoList = document.getElementById('info-list');

    const header = document.createElement('h2');
    header.textContent = titleText;
    header.style.textAlign = 'center'; // Выравнивание по центру
    header.style.marginTop = '20px'; // Отступ сверху 20 пикселей
    header.style.borderBottom = '2px solid #ccc'; // Подчеркивание

    infoList.appendChild(header);
}



export function createPagination(placeholderName,prevPage, currentPage, nextPage, btnClickHandler) {
    const btnPlace = document.getElementById(placeholderName);
    btnPlace.innerHTML = ''; // Очищаем контейнер от предыдущей пагинации


    // Создаем кнопки пагинации
    let prevBtn = document.createElement('button');
    prevBtn.classList.add('btn', 'btn-secondary', 'mx-2');
    prevBtn.textContent = 'Предыдущая';
    prevBtn.addEventListener('click', function (){
        btnClickHandler(prevBtn, localStorage.getItem('token'))
    })
    prevBtn.disabled = prevPage === currentPage; // Отключаем кнопку, если предыдущая страница === текущей

    let nextBtn = document.createElement('button');
    nextBtn.classList.add('btn', 'btn-secondary', 'mx-2');
    nextBtn.textContent = 'Следующая';
    nextBtn.addEventListener('click', function (){
        btnClickHandler(nextPage, localStorage.getItem('token'))
    })
    nextBtn.disabled = nextPage === currentPage; // Отключаем кнопку, если следующая страница === текущей

    // Создаем элемент для отображения текущей страницы
    let currentPageDiv = document.createElement('div');
    currentPageDiv.classList.add('btn', 'btn-primary', 'mx-2');
    currentPageDiv.textContent = currentPage+1

    // Добавляем кнопки в контейнер
    btnPlace.appendChild(prevBtn);
    btnPlace.appendChild(currentPageDiv);
    btnPlace.appendChild(nextBtn);


}