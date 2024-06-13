import * as modalFunc from './pakages/modal.js';
import * as lessonsFunctions from './pakages/lessons.js';
    import * as tokenFunctions from './pakages/token.js';


    let cur_page=0;
    const limit = 10;
    const activitiesList = document.getElementById("activities-list");
    const profileElement = document.getElementById("profile-details");



    // Обработчик события клика на элемент пагинации
    document.querySelectorAll('.pagination .page-item').forEach(item => {
        item.addEventListener('click', function(event) {
            // Получаем значение дня из атрибута data-day
            const day = parseInt(event.target.closest('.page-item').dataset.day);

            // Очистка списка занятий
            activitiesList.style.opacity = 0; // Скрыть список
            // Выполните действия, связанные с выбранным днем
            loadLessons(day);
        });
    });
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
                    // Если ответ 200, значит токен действителен
                    return response.json();
                } else {
                    tokenFunctions.removeToken()
                    // Если ответ не 200, переходим на блок else
                    throw new Error('Токен недействителен');
                }
            }).then(data => {
                // Получение имени и электронной почты из данных токена
                // Получение имени и электронной почты из данных токена
                const tokenValues = tokenFunctions.getTokenValues(token)
                const userName = "Имя: " + tokenValues.name;
                const userEmail = "Логин: " + tokenValues.email;

                var userNameElement = document.createElement("span");
                userNameElement.textContent = userName;

                var userEmailElement = document.createElement("span");
                userEmailElement.textContent = userEmail;

                // Создание кнопки "Выйти" с классами Bootstrap
                var logoutButton = document.createElement("button");
                logoutButton.textContent = "Выйти";
                logoutButton.classList.add("btn", "btn-primary", "mr-2"); // Добавляем классы Bootstrap
                logoutButton.addEventListener("click", function() {
                    // Удаление токена из localStorage
                    tokenFunctions.removeToken()
                    // Перезагрузка страницы
                    location.reload();
                });

                // Добавление элементов на страницу
                profileElement.appendChild(userEmailElement);
                profileElement.appendChild(userNameElement);
                profileElement.appendChild(logoutButton);
            })
            .catch(error => {
                // Блок else: токен недействителен или произошла ошибка при запросе
                console.error('Ошибка при проверке токена:', error);

                // Создание кнопки "Войти"
                var loginButton = document.createElement("button");
                loginButton.textContent = "Войти";
                loginButton.classList.add("btn", "btn-primary"); // Добавляем классы Bootstrap
                // Установка обработчика событий для кнопки
                loginButton.addEventListener("click", function() {
                    window.location.href = "login-register.html";
                });
                // Добавление кнопки на страницу
                profileElement.appendChild(loginButton);
            });

        // Запрос списка занятий
        loadLessons(cur_page)
    });

function loadLessons(day){
    fetch('http://localhost:8080/a?weekDay=' + encodeURIComponent(day), {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            // Проверка успешности ответа
            if (!response.ok) {
                throw new Error('Ошибка сети: ' + response.statusText);
            }
            // Преобразование ответа в формат JSON
            return response.json();
        })
        .then(data => {
            setTimeout(() => {
                activitiesList.innerHTML = ''; // Очищаем список
                activitiesList.style.opacity = 1; // Показать список

                // Добавьте элементы в список
                data.list.forEach(lesson => {
                    addLessonToList(lesson);
                });

            }, 500);
            // Добавьте код для отображения данных на странице
        })
        .catch(error => {
            // Обработка ошибок
            console.error('Ошибка при получении данных:', error);
        });
}

    // Функция для добавления занятия в список
    function addLessonToList(lesson) {
        // Создание элемента списка
        let listItem = document.createElement("li");
        // Добавление класса для оформления
        listItem.classList.add("activities-list-item");
        let lessonId = lesson.LessonId
        // Создание текста с параметрами занятия
        // Добавление текста в элемент списка
        listItem.innerHTML = "ID: " + lessonId + "<br>" +
            "Название: " + lesson.Title + "<br>" +
            "Время начала занятия: " + lessonsFunctions.getTimeString(lesson.Time) + "<br>" +
            "ID тренера: " + lesson.TrainerId + "<br>" +
            "Количество свободных мест: " + lesson.FreeSeats + "/"+ lesson.AvailableSeats + "<br>" +
            "Cложность занятия: " + lesson.Difficult;
        // Создание кнопки "Записаться" с классами Bootstrap
        var signUpButton = document.createElement("button");
        signUpButton.textContent = "Записаться";
        signUpButton.classList.add("btn", "btn-primary", "ml-2"); // Добавляем классы Bootstrap
        signUpButton.style.display = "block";
        // Добавление обработчика события для кнопки

        signUpButton.addEventListener("click", async function() {
            const msg = await lessonsFunctions.handleSubscribeButton(lessonId);
            modalFunc.showInfoModalWithMessage(msg);
        });

        // Добавление кнопки "Записаться" к элементу списка
        listItem.appendChild(signUpButton);

        // Добавление элемента списка в список занятий
        activitiesList.appendChild(listItem);
    }


