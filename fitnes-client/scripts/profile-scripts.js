import * as tokenFunctions from './pakages/token.js';
import * as profileFunctions from './profile/profile-func.js';
import * as lessonsFunctions from './pakages/lessons.js';
import * as modalFunctions from './profile/modal.js';
import * as allModal from './pakages/modal.js';
import {getTimeString, getWeekDayString} from "./pakages/lessons.js";
import {getJwtToken} from "./pakages/token.js";




const roleAdmin = 'Admin'
const roleUser = "User"
const roleTrainer = "Trainer"
const roleNew = "New"
let curPage = 0;
const limit = 5;
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
                userId = tokenPayload.uid
                userToken = token
                profileFunctions.loadUserProfile(token, tokenPayload)
                const userRole = tokenPayload.role
                console.log(userRole)
                // Токен действителен, добавляем элементы на страницу
                switch (userRole) {
                    case roleAdmin:
                        renderAdminMenu(curPage,token);
                        break;
                    case roleTrainer:
                        renderTrainerMenu(curPage,token);
                        break;
                    case roleUser:
                        renderUserMenu(token);
                        break;
                    case roleNew:
                        renderNewUserMenu(token);
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
function renderNewUserMenu(token){
    profileFunctions.addTitle("Пока что тут пусто, дождитесь пока ваше членство подтвердит администратор.")
}



// ======================================================================================
let next_page;
let prev_page;


function renderAdminMenu(page,token) {
    // Реализация для меню администратора
    // Например, добавление пунктов меню для управления пользователями и занятиями
    profileFunctions.addHeader('Список пользователей')
    drawUserList(page, token)


}
function drawUserList(page,token){
    getUserFromServer(page, token)
        .then((data) => {
            if (data) {
                next_page = data.next_page
                prev_page = data.pre_page
                curPage = data.cur_page
                profileFunctions.createPagination('low-btn-place',prev_page,curPage,next_page, drawUserList)
                profileFunctions.createPagination('btn-place',prev_page,curPage,next_page,drawUserList)
                const infoList = document.getElementById("info-list");
                infoList.innerHTML = ''
                data.list.forEach(user =>{
                    addUserToPage(user)
                })
            } else {
                throw new Error();
            }
        })
        .catch((error) => {
            allModal.showInfoModalWithMessage("Ошибка при получении данных")
        });
}

function getUserFromServer(page, token){
    return fetch(
        `http://localhost:8080/account/for-admin/users?page=${encodeURIComponent(
            page
        )}&limit=${encodeURIComponent(limit)}`,
        {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        }
    )
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Ошибка сети: ${response.statusText}`);
            }
            return response.json();
        })
        .catch((error) => {
            console.error('Ошибка при получении данных:', error);
            return null;
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
        "Роль: " + profileFunctions.getRussianRoleName(user.role) + "<br>";


    const editProfileButton = profileFunctions.createButtonWithTextAndHandler('Изменить данные',function (){
        modalFunctions.showModalUserDataEdit(user.id, user.name,user.surname, user.patronymic)
    })
    listItemUser.appendChild(editProfileButton);

    const editRoleButton = profileFunctions.createButtonWithTextAndHandler('Изменить роль',function(){
        modalFunctions.showEditUserRoleModal(user.role, user.id)
    })
    listItemUser.appendChild(editRoleButton);

    const editPasswordButton = profileFunctions.createButtonWithTextAndHandler('Изменить пароль', function (){
        modalFunctions.showChangePasswordModal(user.id)
    })
    listItemUser.appendChild(editPasswordButton);

    // Добавление элемента списка в список занятий
    document.getElementById("info-list").appendChild(listItemUser);
}


// ======================================================================================
function renderTrainerMenu(page, token){
    profileFunctions.addHeader('Проводимые занятия')
    const createLessonBtn = profileFunctions.createButtonWithTextAndHandler('Создать занятие',handleCreateButtonClick)
    createLessonBtn.style.marginBottom = '20px'
    document.getElementById('create-btn-place').appendChild(createLessonBtn)

    drawTrainerLessonsList(page, token)
}
function drawTrainerLessonsList(page, token){
    // Отрисовка
    getTrainerLessonsFromServer(page, token)
        .then((data) => {
            if (data) {
                next_page = data.next_page
                prev_page = data.pre_page
                curPage = data.cur_page
                profileFunctions.createPagination('low-btn-place',prev_page,curPage,next_page, drawTrainerLessonsList)
                profileFunctions.createPagination('btn-place',prev_page,curPage,next_page,drawTrainerLessonsList)
                const infoList = document.getElementById("info-list");
                infoList.innerHTML = ''
                if (data.list.length === 0){
                    profileFunctions.addTitle('Тут пока пусто')
                }else{
                    data.list.forEach(user =>{
                        addTrainerLessonsToList(user)
                    })
                }
            } else {
                throw new Error();
            }
        })
        .catch((error) => {
            allModal.showInfoModalWithMessage("Ошибка при получении данных")
        });
}
function getTrainerLessonsFromServer(page, token) {
    return fetch('http://localhost:8080/lessons/get/by-trainer?page=' + encodeURIComponent(page) + '&limit=' + encodeURIComponent(limit) + '&trainerId=' + encodeURIComponent(userId), {
        method: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Ошибка сети: ${response.statusText}`);
            }
            return response.json();
        })
        .catch((error) => {
            console.error('Ошибка при получении данных:', error);
            return null;
        });
}


function addTrainerLessonsToList(lesson){
    // Создание элемента списка
    let listItem = document.createElement("li");
    // Добавление класса для оформления
    listItem.classList.add("activities-list-item");

    // Создание текста с параметрами занятия
    // Добавление текста в элемент списка
    const startTime = lessonsFunctions.getTimeString(lesson.Time)
    const weekDay = getWeekDayString(lesson.DayOfWeek)
    listItem.innerHTML = "ID: " + lesson.LessonId + "<br>" +
        "Название: " + lesson.Title + "<br>" +
        "Время начала: " + startTime + "<br>" +
        "Количество свободных мест: " + lesson.FreeSeats + "/"+ lesson.AvailableSeats + "<br>" +
        "Сложность: " + lesson.Difficult + "<br>"+
        "День недели: "+ weekDay;

    // Создание кнопки "Записаться" с классами Bootstrap


    const signUpButton = profileFunctions.createButtonWithTextAndHandler('Удалить занятие',function (){
        modalFunctions.showDeleteModal(lesson.LessonId, lesson.Title, weekDay, startTime);
    })
    const editBtn = profileFunctions.createButtonWithTextAndHandler('Редактировать занятие',function (){
        modalFunctions.showEditLessonModal(userId,lesson.LessonId ,lesson.Title, lesson.AvailableSeats,lesson.Description, lesson.Difficult,lesson.DayOfWeek,startTime);
        drawTrainerLessonsList(curPage, localStorage.getItem('token'))
    })
    // Добавление кнопки "Записаться" к элементу списка
    listItem.appendChild(signUpButton);
    listItem.appendChild(editBtn)

    // Добавление элемента списка в список занятий
    document.getElementById("info-list").appendChild(listItem);

}

function handleCreateButtonClick() {
    modalFunctions.showCreateLessonModal(userId)
}

// ======================================================================================
function renderUserMenu(token){
    profileFunctions.addHeader('Мои занятия')
    drawUserLessons(curPage, token)
}
function drawUserLessons(page, token){
    // Отрисовка
    getUserLessonsFromServer(page, token)
        .then((data) => {
            if (data) {
                next_page = data.next_page
                prev_page = data.pre_page
                curPage = data.cur_page
                profileFunctions.createPagination('low-btn-place',prev_page,curPage,next_page, drawUserLessons)
                profileFunctions.createPagination('btn-place',prev_page,curPage,next_page,drawUserLessons)
                const infoList = document.getElementById("info-list");
                infoList.innerHTML = ''
                if (data.list.length === 0){
                    profileFunctions.addTitle('Тут пока пусто')
                }else{
                    data.list.forEach(user =>{
                        addUserLessonToList(user)
                    })
                }
            } else {
                throw new Error();
            }
        })
        .catch((error) => {
            console.log(error)
            allModal.showInfoModalWithMessage("Ошибка при получении данных")
        });
}
function getUserLessonsFromServer(page, token){
    return fetch('http://localhost:8080/lessons/get/by-user?page=' + encodeURIComponent(page) + '&limit=' + encodeURIComponent(limit) + '&trainerId=' + encodeURIComponent(userId), {
        method: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Ошибка сети: ${response.statusText}`);
            }
            return response.json();
        })
        .catch((error) => {
            console.error('Ошибка при получении данных:', error);
            return null;
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
    const startTime = lessonsFunctions.getTimeString(lesson.Time)
    const weekDay = getWeekDayString(lesson.DayOfWeek)
    listItem.innerHTML = "ID: " + lesson.LessonId + "<br>" +
        "Название: " + lesson.Title + "<br>" +
        "Время начала: " + startTime + "<br>" +
        "Количество свободных мест: " + lesson.FreeSeats + "/"+ lesson.AvailableSeats + "<br>" +
        "Сложность: " + lesson.Difficult + "<br>"+
        "День недели: "+ weekDay;

    // Создание кнопки "Записаться" с классами Bootstrap


    const signUpButton = profileFunctions.createButtonWithTextAndHandler('Отписаться от занятия',function (){
        lessonsFunctions.handleSubscribeButton(lesson.LessonId);
        drawUserLessons(curPage, localStorage.getItem('token'))
    })
    // Добавление кнопки "Записаться" к элементу списка
    listItem.appendChild(signUpButton);

    // Добавление элемента списка в список занятий
    document.getElementById("info-list").appendChild(listItem);

}


// ======================================================================================
// Функция для отображения занятий пользователя на странице
/*
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
*/

// =========================================================================
// // изменение своей страницы
// function setupLogout(){
//     location.reload()
// }
// Получаем кнопку по ее ID
const logoutButton = document.getElementById('logoutButton');

// Добавляем обработчик события click
logoutButton.addEventListener('click', function() {
    // Вызываем функцию exitFromAccount() при клике на кнопку
    modalFunctions.showExitModal()
});


const editMyProfileButton = document.getElementById('editProfileButton')
editMyProfileButton.addEventListener('click',function (){
    const userData = profileFunctions.parseUserData()
    modalFunctions.showModalUserDataEdit(userData.id, userData.name, userData.surname, userData.patronymic)
})

const editMyPasswordButton = document.getElementById("editPasswordButton")
editMyPasswordButton.addEventListener('click',function (){
    modalFunctions.showChangePasswordModal(userId)
})

// =============================================================================