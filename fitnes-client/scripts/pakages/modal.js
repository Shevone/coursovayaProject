function openInfoModal(title, bodyContent, footerContent, handler) {
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
export function showInfoModalWithMessage(message, handler) {
    // Функция для открытия модального окна с подтверждением
    const title = "Сообщение";
    const bodyContent = `<p>${message}</p>`; // Добавляем сообщение в bodyContent
    const footerContent = `
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" id="closeBtn">Ок</button>
    `;
    openInfoModal(title, bodyContent, footerContent, handler);
}