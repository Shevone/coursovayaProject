import * as tokenFunctions from './pakages/token.js';
import * as profileFunctions from './profile/profile-func.js';
import * as lessonsFunctions from './pakages/lessons.js';
import * as modalFunctions from './profile/modal.js';


const roleAdmin = 'Admin'
const roleUser = "User"
const roleTrainer = "Trainer"
let curPage = 0;
const limit = 10;
let userId = -1;
let userToken = "";

document.addEventListener("DOMContentLoaded", function() {
    let token = tokenFunctions.getJwtToken()

    fetch('http://localhost:8080/account/token-valid', {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    })
        .then(response => {
            if (response.status === 200) {
                const tokenPayload = tokenFunctions.getTokenPayload(token)
                userId = tokenPayload.id
                userToken = token
                profileFunctions.loadUserProfile(token, tokenPayload)
                const userRole = tokenFunctions.getUserRoleFromToken(token);
                console.log(userRole)
                // Токен действителен, добавляем элементы на страницу
                switch (userRole) {
                    case roleAdmin:
                        renderAdminMenu(token);
                        break;
                    case roleTrainer:
                        renderTrainerMenu(token);
                        break;
                    case roleUser:
                        renderUserMenu(token);
                        break;
                    default:
                        // Если роль неизвестна или не определена, выводим сообщение об ошибке
                        console.error('Не удалось определить роль пользователя.');
                        break;
                }
            } else {
                // Если ответ не 200, токен недействителен
                modalFunctions.showConfirmationModal();
            }
        })
});

// ======================================================================================
function renderAdminMenu(token) {
    // Реализация для меню администратора
    // Например, добавление пунктов меню для управления пользователями и занятиями
    profileFunctions.addHeader('Список пользователей')
    fetch('http://localhost:8080/account/for-admin/users?page='+ encodeURIComponent(curPage) + '&limit='+ encodeURIComponent(limit), {
        method: 'GET',
        headers: {
            'Authorization': "Bearer " + token
        }
    }).then(response => {
        // Проверка успешности ответа
        if (!response.ok) {
            throw new Error('Ошибка сети: ' + response.statusText);
        }
        // Преобразование ответа в формат JSON
        return response.json();
    }).then(data => {
        // Обработка данных занятий
        console.log(data);
        data.list.forEach(
            user => addUserToPage(user)
        )
        // Добавьте код для отображения данных на странице
    })
        .catch(error => {
            // Обработка ошибок
            console.error('Ошибка при получении данных:', error);
        });


}

function addUserToPage(user){
    if (user.id === userId ){
        return
    }
    // Создание элемента списка
    let listItemUser = document.createElement("li");
    // Добавление класса для оформления
    listItemUser.classList.add("activities-list-item");

    // Создание текста с параметрами занятия
    // Добавление текста в элемент списка
    listItemUser.innerHTML = "Логин: " + user.email + "<br>" +
        "Id: " + user.id + "<br>" +
        "Имя: " + user.name + "<br>" +
        "Фамилия: " + user.surname + "<br>" +
        "Отчество: " + user.patronymic + "<br>" +
        "Номер телефона: " + user.phoneNumber + "<br>" +
        "Роль: " + user.role + "<br>";

    // TODO кнопки: функция обработки каждой кнопки
    const editProfileButton = profileFunctions.createButtonWithTextAndHandler('Изменить данные',modalFunctions.showModalUserDataEdit(user.id, user.name,user.surname, user.patronymic))
    listItemUser.appendChild(editProfileButton);

    const editRoleButton = profileFunctions.createButtonWithTextAndHandler('Изменить роль',location.reload())
    listItemUser.appendChild(editRoleButton);

    const editPasswordButton = profileFunctions.createButtonWithTextAndHandler('Изменить пароль',location.reload())
    listItemUser.appendChild(editPasswordButton);

    // Добавление элемента списка в список занятий
    document.getElementById("info-list").appendChild(listItemUser);
}

// ======================================================================================
function editUserRole(userId, curRole){

}
// ======================================================================================
function editUserPassword(userId){

}
// ======================================================================================
function renderTrainerMenu() {
    // Реализация для меню тренера
    // Например, добавление пунктов меню для управления занятиями и просмотра профиля
    profileFunctions.addHeader('Проводимые занятия')
}

// ======================================================================================
function renderUserMenu(token) {
    profileFunctions.addHeader('Мои занятия')
    fetch('http://localhost:8080/lessons/get/by-user?page='+ encodeURIComponent(curPage) + '&limit='+ encodeURIComponent(limit), {
        method: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    }).then(response => {
        // Проверка успешности ответа
        if (!response.ok) {
            throw new Error('Ошибка сети: ' + response.statusText);
        }
        // Преобразование ответа в формат JSON
        return response.json();
    }).then(data => {
        // Обработка данных занятий
        console.log(data);
        data.list.forEach(lesson =>
            addUserLessonToList(lesson)
        )
        // Добавьте код для отображения данных на странице
    })
        .catch(error => {
            // Обработка ошибок
            console.error('Ошибка при получении данных:', error);
        });
}


// Функция для добавления занятия в список
function addUserLessonToList(lesson) {
    // Создание элемента списка
    let listItem = document.createElement("li");
    // Добавление класса для оформления
    listItem.classList.add("activities-list-item");

    // Создание текста с параметрами занятия
    // Добавление текста в элемент списка
    listItem.innerHTML = "ID: " + lesson.id + "<br>" +
        "Title: " + lesson.title + "<br>" +
        "Time: " + lesson.time + "<br>" +
        "Trainer ID: " + lesson.trainer_id + "<br>" +
        "Available Seats: " + lesson.available_seats + "<br>" +
        "Difficult: " + lesson.difficult;

    // Создание кнопки "Записаться" с классами Bootstrap


    const signUpButton = profileFunctions.createButtonWithTextAndHandler('Удалить(отписаться)',function (){
        lessonsFunctions.handleSubscribeButton(lesson.id);
        location.reload();
    })
    // Добавление кнопки "Записаться" к элементу списка
    listItem.appendChild(signUpButton);

    // Добавление элемента списка в список занятий
    document.getElementById("info-list").appendChild(listItem);
}


// ======================================================================================
// Функция для отображения занятий пользователя на странице
function displayUserLessons(userLessons) {
    var userLessonsList = document.getElementById("user-lessons-list");

    // Очистка списка перед добавлением новых данных
    userLessonsList.innerHTML = '';

    // Перебор полученных занятий и добавление их в список
    userLessons.forEach(function(lesson) {
        var listItem = document.createElement("li");
        listItem.classList.add("list-group-item");
        listItem.textContent = "Название занятия: " + lesson.title + ", Дата: " + lesson.date + ", Время: " + lesson.time;
        userLessonsList.appendChild(listItem);
    });
}

// =========================================================================
// изменение своей страницы
function setupLogout(){
    location.reload()
}
function exitFromAccount(){
    modalFunctions.showExitModal()
}