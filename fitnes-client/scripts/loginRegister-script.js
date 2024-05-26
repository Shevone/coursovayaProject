

document.addEventListener("DOMContentLoaded", function() {
    // Проверяем наличие токена в localStorage
    const token = localStorage.getItem('token')
    fetch("http://localhost:8080/account/token-valid", {
        method: "POST",
        headers: {
            "Authorization": "Bearer " + token
        }
    })
        .then(response => {
            // Если ответ 200, перенаправляем на index.html
            if (response.status === 200) {
                window.location.href = "index.html";
            } else {
                // Если ответ не 200, удаляем токен из localStorage
                localStorage.removeItem("token");
            }
        })
        .catch(error => {
            console.error("Error:", error);
            // В случае ошибки также удаляем токен из localStorage
            localStorage.removeItem("token");
        });
    document.getElementById('openRegModal').addEventListener('click', function() {
        showRegisterModal()
    });
    const btn =  document.getElementById('loginBtn')
   btn.addEventListener('click', function (){
        login()
    })
});


document.getElementById("closePopup").addEventListener("click", function() {
    document.getElementById("overlay").style.display = "none";
    document.getElementById("registerPopup").style.display = "none";
});

function login (){

    let username = document.getElementById("loginUsername").value;
    let password = document.getElementById("loginPassword").value;
    fetch('http://localhost:8080/account/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(
            {
                email : username,
                password: password
            }
        ),
    })
        .then(response => {
            if (response.ok) {
                // Преобразуем тело ответа в формат JSON
                return response.json();
            } else {
                // Если ответ не успешен, бросаем ошибку
                throw new Error('login failed!');
            }
        })
        .then(data => {
            // Получаем токен из объекта данных
            let token = data.token;
            console.log("Token:", token);

            // Далее вы можете использовать этот токен для аутентификации и авторизации пользователя
            alert("login successful!");
             if (token) {
                // Сохраняем токен в локальном хранилище (localStorage) или куках, чтобы его можно было использовать на других страницах
                localStorage.setItem('token', token);

                // Перенаправляем пользователя на другую страницу
                 window.location.href = "index.html";
            } else {
                // Если токен не был получен, обрабатываем ошибку аутентификации
                console.error("Ошибка аутентификации: токен не был получен");
                alert("Войдите еще раз")
            }
        })
        .catch(error => {
            // Обрабатываем ошибку, если что-то пошло не так
            alert(error.message);
        });

}

function showRegisterModal(){
    const title = "Регистрация";
    const bodyContent = `
        <form id="registerForm">
            <input type="text" id="Name" placeholder="Имя" required><br>
            <input type="text" id="Surname" placeholder="Фамилия" required><br>
            <input type="text" id="Patronymic" placeholder="Отчество(при наличии)" required><br>
            <input type="text" id="phoneNumber" placeholder="Номер телефона" required><br>
            <input type="text" id="email" placeholder="Логин" required><br>
            <input type="password" id="registerPassword" placeholder="Пароль" required><br>
            <input type="password" id="confirmPassword" placeholder="Подтвердите пароль" required><br>
        </form>
    `;
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
        <button type="button" class="btn btn-primary" id="confirmButton">Зарегестрироваться!</button>
    `;
    openDynamicModal(title, bodyContent, footerContent, registerFetch);
}
function registerFetch(){
    let name = document.getElementById("Name").value;
    let surname = document.getElementById("Surname").value;
    let patronymic = document.getElementById("Patronymic").value;
    let phoneNumber = document.getElementById("phoneNumber").value;
    let email = document.getElementById("email").value;
    let password = document.getElementById("registerPassword").value;
    let confirmPassword = document.getElementById("confirmPassword").value;

    if (password !== confirmPassword) {
        alert("Passwords do not match!");
        return;
    }


    let userData = {
        name: name,
        surname: surname,
        patronymic: patronymic,
        phoneNumber: phoneNumber,
        email: email,
        password: password
    };
    fetch('http://localhost:8080/account/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    }).then(response =>{
        if (response.ok){
            alert("Регистрация прошла успешно!")
            location.reload()
        }else {
            alert("Ошибка при регистрации")
        }
    })
}

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
