export function showRegisterModal(){
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
