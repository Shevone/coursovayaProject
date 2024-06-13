function openDynamicModal(title, bodyContent, footerContent, handlerMethod) {
    // Заполняем заголовок модального окна
    document.getElementById('dynamicModalLabel').innerText = title;

    // Заполняем содержимое тела модального окна
    document.getElementById('dynamicModalBody').innerHTML = bodyContent;

    // Заполняем содержимое подвала модального окна
    document.getElementById('dynamicModalFooter').innerHTML = footerContent;

    // Открываем модальное окно
    let modal = new bootstrap.Modal(document.getElementById('dynamicModal'));
    modal.show();
    const confirmLogoutButton = document.getElementById('confirmButton');
    confirmLogoutButton.addEventListener('click', handlerMethod);
}

function openInfoModalProfile(title, bodyContent, footerContent, handler) {
    // Заполняем заголовок модального окна
    document.getElementById('dynamicModalInfoLabel').innerText = title;

    // Заполняем содержимое тела модального окна
    document.getElementById('dynamicModalInfoBody').innerHTML = bodyContent;

    // Заполняем содержимое подвала модального окна
    document.getElementById('dynamicModalInfoFooter').innerHTML = footerContent;

    // Открываем модальное окно
    let modal = new bootstrap.Modal(document.getElementById('dynamicModalInfo'));
    modal.show();
    const closeButton = document.getElementById('closeBtn');
    closeButton.addEventListener('click', handler)

}


// ==================================================================================================
function showInfoModalWithMessage(message, handler) {
    // Функция для открытия модального окна с подтверждением
    const title = "Сообщение";
    const bodyContent = `<p>${message}</p>`; // Добавляем сообщение в bodyContent
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" id="closeBtn">Ок</button>
    `;
    if (handler == null){
        handler = function (){}
    }
    openInfoModalProfile(title, bodyContent, footerContent, handler);
}
// ==================================================================================================
export function showConfirmationModal() {
    // Функция для открытия модального окна с подтверждением
    const title = "Подтверждение";
    const bodyContent = "<p>Для доступа к личному кабинету необходимо авторизоваться.</p>";
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton"">Войти</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, redirectToLoginPage);
}
function redirectToLoginPage() {
    // Перенаправление на страницу входа
    window.location.href = "login-register.html";
}


// ==================================================================================================
export function showModalUserDataEdit(userId, curName, curSurname, curPatronymic){
    const title = "Редактирование профиля";
    const bodyContent = `
        <form id="editProfileForm">
            <div class="mb-3">
                <label for="userIdInput" class="form-label">ID пользователя</label>
                 <input type="text" class="form-control" id="userIdInput" placeholder="${userId}" readonly>
            </div>
            <div class="mb-3">
                <label for="firstNameInput" class="form-label">Имя</label>
                <input type="text" class="form-control" id="firstNameInput" placeholder="${curName}">
            </div>
            <div class="mb-3">
                <label for="lastNameInput" class="form-label">Фамилия</label>
                <input type="text" class="form-control" id="lastNameInput" placeholder="${curSurname}">
            </div>
            <div class="mb-3">
                <label for="patronymicInput" class="form-label">Отчество</label>
                <input type="text" class="form-control" id="patronymicInput" placeholder="${curPatronymic}">
            </div>
        </form>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton">Обновить данные</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, updateProfile);
}

function updateProfile(event){
    // Получаем значения из полей формы
    let firstName = document.getElementById("firstNameInput").value.trim();
    let lastName = document.getElementById("lastNameInput").value.trim();
    let patronymic = document.getElementById("patronymicInput").value.trim();
    let userId = document.getElementById("userIdInput").placeholder;

    // Проверяем, что имя и фамилия не пустые
    if (firstName === "" || lastName === "") {
        showInfoModalWithMessage("Имя и фамилия должны быть заполнены.");
        return;
    }

    // Формируем данные для запроса
    let requestData = {
        userId, // Предполагается, что у вас есть функция getUserId(), которая возвращает ID пользователя
        name: firstName,
        surname: lastName,
        patronymic: patronymic
    };
    const userToken = localStorage.getItem('token')
    fetch('http://localhost:8080/account/edit-profile', {
        method: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + userToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
        .then(response => response.json())
        .then(data => {
            showInfoModalWithMessage(data.message, function (){
                location.reload()
            });


        })
        .catch(error => {
            showInfoModalWithMessage(error, function (){

                location.reload()
            }); // Выводим результат запроса пользователю


        });
    event.target.removeEventListener('click', updateProfile);
}

// ==================================================================================================
export function showExitModal(){
    const title = "Подтверждение выхода";
    const bodyContent = `
        <div class="modal-body">
                Вы уверены, что хотите выйти из аккаунта?
        </div>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton">Да, выйти</button>
    `;
    openDynamicModal(title, bodyContent, footerContent,confirmLogout);
}
function confirmLogout(event) {
    localStorage.removeItem("token");
    window.location.href = "index.html";
    event.target.removeEventListener('click', confirmLogout); // Удаление обработчика
}
// ==================================================================================================

const roleAdmin = 'Admin';
const roleUser = 'User';
const roleTrainer = 'Trainer';
const roleNew = 'New'

export function showEditUserRoleModal(currentRole, userId) {
    const title = "Изменение роли пользователя";
    const bodyContent = `
        <div class="modal-body">
            <div class="form-group">
                <label for="userId">ID пользователя:</label>
                <input type="text" class="form-control" id="userId" value="${userId}" readonly> 
            </div>
            <div class="form-group">
                <label for="userRoleSelect">Выберите роль:</label>
                <select class="form-control" id="userRoleSelect">
                    <option value=${roleNew} ${currentRole === roleNew ? 'selected' : ''}>Гость</option>
                    <option value=${roleAdmin} ${currentRole === roleAdmin ? 'selected' : ''}>Админ</option>
                    <option value=${roleUser} ${currentRole === roleUser ? 'selected' : ''}>Клиент</option>
                    <option value=${roleTrainer} ${currentRole === roleTrainer ? 'selected' : ''}>Тренер</option>
                </select>
            </div>
        </div>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton">Сохранить</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, updateUserRole);
}
function updateUserRole(event){
    const userId = document.getElementById('userId').value;
    const selectedRole = document.getElementById('userRoleSelect').value;
    let selectedRoleInt
    if (selectedRole === roleUser){
        selectedRoleInt = 0
    }else if (selectedRole === roleTrainer){
        selectedRoleInt = 1
    }else if (selectedRole === roleAdmin) {
        selectedRoleInt = 2
    }else {
        selectedRoleInt = 3
    }

    let requestData = {
        user_id: parseInt(userId),
        new_role: selectedRoleInt,

    };

    const userToken = localStorage.getItem('token')
    fetch('http://localhost:8080/account/for-admin/update-role', {
        method: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + userToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    }).then(response => response.json())
        .then(data => {

            showInfoModalWithMessage(data, function (){
                event.target.removeEventListener('click', updateUserRole);
                location.reload()
            }); // Выводим результат запроса пользователю


        })
        .catch(error => {
            showInfoModalWithMessage(error, function (){
                event.target.removeEventListener('click', updateUserRole);
                location.reload()
            }); // Выводим результат запроса пользователю


        });
}
// ==================================================================================================
export function showChangePasswordModal(userId) {
    const title = "Изменение пароля";
    const bodyContent = `
        <div class="modal-body">
            <div class="form-group">
                <label for="userId">ID пользователя:</label>
                <input type="text" class="form-control" id="userId" value="${userId}" readonly> 
            </div>
            <div class="form-group">
                <label for="newPassword">Новый пароль:</label>
                <input type="password" class="form-control" id="newPassword" required>
            </div>
            <div class="form-group">
                <label for="confirmPassword">Подтверждение пароля:</label>
                <input type="password" class="form-control" id="confirmPassword" required>
            </div>
        </div>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton">Сохранить</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, changePassword);
}
async function changePassword(event) {
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const userId = parseInt(document.getElementById('userId').value);

    if (newPassword !== confirmPassword) {
        showInfoModalWithMessage("Пароли не совпадают!");
        return;
    }

    const token = localStorage.getItem('token');
    const request ={
        user_id: userId,
        password: newPassword
    }
    try{
        const response = await fetch('http://localhost:8080/account/edit-password', {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(request)
        });
        // Проверяем, был ли ответ успешным
        if (!response.ok) {
            throw new Error('Ошибка сети: ' + response.statusText);
        }

        // Преобразуем ответ в формат JSON
        const data = await response.json();
        if (response.ok) {
            const msg = data;
            showInfoModalWithMessage(msg, function (){location.reload()}); // Вывод сообщения из ответа сервера
        }
    } catch (error){
        console.log(error)
        showInfoModalWithMessage('Ошибка при изменении пароля', function (){})
    }


}
// ==================================================================================================
export function showCreateLessonModal(trainerId) {
    const title = "Создать занятие";
    const bodyContent = `
    <div class="modal-body">
      <div class="form-group">
        <label for="userId">ID тренера:</label>
        <input type="text" class="form-control" id="userId" value="${trainerId}" readonly> 
      </div>
      <div class="form-group">
        <label for="lessonName">Название занятия:</label>
        <input type="text" class="form-control" id="lessonName" required>
      </div>
      <div class="form-group">
        <label for="availableSeats">Количество свободных мест:</label>
        <div class="input-group">
          <button class="btn btn-outline-secondary" id="decrementSeats">-</button>
          <input type="number" class="form-control" id="availableSeats" value="1" min="1" required>
          <button class="btn btn-outline-secondary" id="incrementSeats">+</button>
        </div>
      </div>
      <div class="form-group">
        <label for="lessonDescription">Описание:</label>
        <textarea class="form-control" id="lessonDescription" rows="3" required></textarea>
      </div>
      <div class="form-group">
        <label for="lessonDifficulty">Сложность:</label>
        <select class="form-control" id="lessonDifficulty" required>
          <option value="1">Простой</option>
          <option value="MEDIUM">Средний</option>
          <option value="HARD">Сложный</option>
        </select>
      </div>
      <div class="form-group">
         <label for="lessonDifficulty">День недели:</label>
        <select class="form-control" id="lessonWeekDay" required>
          <option value="0">Понедельник</option>
          <option value="1">Вторник</option>
          <option value="2">Среда</option>
          <option value="3">Четверг</option>
          <option value="4">Пятница</option>
          <option value="5">Суббот</option>
          <option value="6">Воскресенье</option>
        </select>
      </div>
      <div class="form-group">
        <label for="lessonTime">Время:</label>
        <input type="time" class="form-control" id="lessonTime" required>
      </div>
    </div>
  `;
    const footerContent = `
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
    <button type="button" class="btn btn-primary" id="confirmButton">Сохранить</button>
  `;
    openDynamicModal(title, bodyContent, footerContent, createLesson);
    handleSeatsChange();
}

function handleSeatsChange() {
    const availableSeatsInput = document.getElementById('availableSeats');
    const decrementSeatsButton = document.getElementById('decrementSeats');
    const incrementSeatsButton = document.getElementById('incrementSeats');

    decrementSeatsButton.addEventListener('click', () => {
        let currentSeats = parseInt(availableSeatsInput.value);
        if (currentSeats > 1) {
            availableSeatsInput.value = currentSeats - 1;
        }
    });

    incrementSeatsButton.addEventListener('click', () => {
        let currentSeats = parseInt(availableSeatsInput.value);
        availableSeatsInput.value = currentSeats + 1;
    });
}

// Добавьте эту строку после создания модального окна

function createLesson() {
    // Получаем данные из формы
    const trainerId = parseInt(document.getElementById('userId').value);
    const lessonName = document.getElementById('lessonName').value;
    const availableSeats = parseInt(document.getElementById('availableSeats').value);
    const lessonDescription = document.getElementById('lessonDescription').value;
    const lessonDifficulty = document.getElementById('lessonDifficulty').value;
    const dayOfWeek = parseInt(document.getElementById('lessonWeekDay').value);
    const lessonTime = document.getElementById('lessonTime').value;

    if (lessonName === ""){
        showInfoModalWithMessage("Название не должно быть пустым")
    }
    if (lessonTime === ""){
        showInfoModalWithMessage("Выберите время занятия")
    }
    const requestData = {
        title : lessonName,
        time : lessonTime,
        dayOfWeek : dayOfWeek,
        trainerId : trainerId,
        availableSeats : availableSeats,
        difficult : lessonDifficulty,
        description : lessonDescription
    }
    fetch('http://localhost:8080/lessons/create', {
        method: 'Post',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token'),
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    }).then(response => {
        if (response.ok) {
            showInfoModalWithMessage("Занятие создано!", function (){location.reload()});

        } else {
            showInfoModalWithMessage("Произошла ошибка, попробуйте позже");
        }
    });
}

export function showDeleteModal(lessonId, lessonTitle, dayOfWeek, lessonStartTime){
    // Функция для открытия модального окна с подтверждением
    const title = "Подтверждение удаления занятия";
    const bodyContent = `
        <p>Вы действительно хотите удалить занятие:</p>
        <input type="hidden" id="lessonIdInput" value="${lessonId}"> 
        <p><strong>Название:</strong> ${lessonTitle}</p>
        <p><strong>День недели:</strong> ${dayOfWeek}</p>
        <p><strong>Время:</strong> ${lessonStartTime}</p>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton" data-lesson-id="${lessonId}">Удалить</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, deleteLesson);
}
function deleteLesson() {

    const lessonIdInput = document.getElementById("lessonIdInput"); // Предполагается, что на странице есть скрытое поле с ID "lessonIdInput"
    const lessonId = parseInt(lessonIdInput.value);
    const requestData = {
        lesson_id: lessonId
    }
    fetch('http://localhost:8080/lessons/delete', {
        method: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token'),
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    }).then(response => {
        if (response.ok) {
            showInfoModalWithMessage("Занятие удалено", function (){ location.reload()});

        } else {
            showInfoModalWithMessage("Произошла ошибка, попробуйте позже")
        }
    });
}


export function showEditLessonModal(trainerId, lessonId,lessonTitle, lessonSeatsCount, lessonDescription, lessonDifficult,lessonDayOfWeek,lessonStartTime){
    const title = "Редактировать занятие";
    const bodyContent = `
    <div class="modal-body">
      <div class="form-group">
        <label for="lessonTrainerId">ID тренера:</label>
        <input type="text" class="form-control" id="lessonTrainerId" value="${trainerId}" readonly> 
      </div>
      <div class="form-group">
        <label for="lessonId">ID занятия:</label>
        <input type="text" class="form-control" id="editLessonId" value="${lessonId}" readonly> 
      </div>
      <div class="form-group">
        <label for="lessonName">Название занятия:</label>
        <input type="text" class="form-control" id="lessonName" value="${lessonTitle}" required>
      </div>
      <div class="form-group">
        <label for="availableSeats">Количество свободных мест:</label>
        <div class="input-group">
          <button class="btn btn-outline-secondary" id="decrementSeats">-</button>
          <input type="number" class="form-control" id="availableSeats" value="${lessonSeatsCount}" min="1" required>
          <button class="btn btn-outline-secondary" id="incrementSeats">+</button>
        </div>
      </div>
      <div class="form-group">
        <label for="lessonDescription">Описание:</label>
        <textarea class="form-control" id="lessonDescription" rows="3" required>${lessonDescription}</textarea>
      </div>
      <div class="form-group">
        <label for="lessonDifficulty">Сложность:</label>
        <select class="form-control" id="lessonDifficulty" required>
          <option value="EASY" ${lessonDifficult === "EASY" ? 'selected' : ''}>Простой</option>
          <option value="MEDIUM" ${lessonDifficult === "MEDIUM" ? 'selected' : ''}>Средний</option>
          <option value="HARD" ${lessonDifficult === "HARD" ? 'selected' : ''}>Сложный</option>
        </select>
      </div>
      <div class="form-group">
         <label for="lessonDifficulty">День недели:</label>
        <select class="form-control" id="lessonWeekDay" required>
          <option value="0" ${lessonDayOfWeek === 0 ? 'selected' : ''}>Понедельник</option>
          <option value="1" ${lessonDayOfWeek === 1 ? 'selected' : ''}>Вторник</option>
          <option value="2" ${lessonDayOfWeek === 2 ? 'selected' : ''}>Среда</option>
          <option value="3" ${lessonDayOfWeek === 3 ? 'selected' : ''}>Четверг</option>
          <option value="4" ${lessonDayOfWeek === 4 ? 'selected' : ''}>Пятница</option>
          <option value="5" ${lessonDayOfWeek === 5 ? 'selected' : ''}>Суббота</option>
          <option value="6" ${lessonDayOfWeek === 6 ? 'selected' : ''}>Воскресенье</option>
        </select>
      </div>
      <div class="form-group">
        <label for="lessonTime">Время:</label>
        <input type="time" class="form-control" id="lessonTime" value="${lessonStartTime}" required>
      </div>
    </div>
  `;
    const footerContent = `
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
    <button type="button" class="btn btn-primary" id="confirmButton">Сохранить</button>
  `;
    openDynamicModal(title, bodyContent, footerContent, editLesson);
    handleSeatsChange();
}
function editLesson() {
    const lessonId = parseInt(document.getElementById('editLessonId').value)
    // Получаем данные из формы
    const trainerId = parseInt(document.getElementById('lessonTrainerId').value);
    const lessonName = document.getElementById('lessonName').value;
    const availableSeats = parseInt(document.getElementById('availableSeats').value);
    const lessonDescription = document.getElementById('lessonDescription').value;
    const lessonDifficulty = document.getElementById('lessonDifficulty').value;
    const dayOfWeek = parseInt(document.getElementById('lessonWeekDay').value)
    const lessonTime = document.getElementById('lessonTime').value;


    let requestData = {

        title: lessonName,
        time: lessonTime,

        difficult: lessonDifficulty,

    }

    const ok =validateRequestData(requestData)
    if (!ok){
        showInfoModalWithMessage("Все данные должны быть заполнены")
        return
    }
    requestData.trainerId = trainerId;
    requestData.lessonId = lessonId;
    requestData.dayOfWeek = dayOfWeek;
    requestData.availableSeats = availableSeats;
    requestData.description = lessonDescription;
    fetch('http://localhost:8080/lessons/edit', {
        method: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token'),
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    }).then(response => {
        if (response.ok) {
            showInfoModalWithMessage("Занятие отредактированно!", function (){ location.reload()});

        } else {
            showInfoModalWithMessage("Произошла ошибка, попробуйте позже")
        }
    });

}
function validateRequestData(requestData) {
    for (const key in requestData) {
        if (requestData[key] === undefined || requestData[key].trim() === "") {
            return false; // Если поле пустое, возвращаем false
        }
    }
    return true; // Если все поля заполнены, возвращаем true
}