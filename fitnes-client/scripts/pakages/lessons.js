

export async function handleSubscribeButton(lessonId) {
    let signUpData = {
        lesson_id: lessonId,
    };
    const token = localStorage.getItem('token');

    if (token) {
        // Если токен есть
        try {
            const response = await fetch('http://localhost:8080/lessons/sign', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify(signUpData),
            });

            // Проверяем, был ли ответ успешным
            if (!response.ok) {
                throw new Error('Ошибка сети: ' + response.statusText);
            }

            // Преобразуем ответ в формат JSON
            const data = await response.json();

            // Возвращаем поле message из тела запроса, если статус код == 200
            if (response.status === 200) {
                const message = data.message;
                return message;
            } else {
                return 'Не удалось записаться';
            }
        } catch (error) {
            // Обработка ошибок
            console.error('Ошибка при получении данных:', error);
            return 'Не удалось записаться';
        }
    } else {
        return 'Необходимо авторизоваться для записи на занятие';
    }
}

export function getWeekDayString(weekDayInt){
    switch (weekDayInt) {
        case 0:
            return 'Понедельник'
        case 1:
            return 'Вторник'
        case 2:
            return 'Среда'
        case 3:
            return 'Четверг'
        case 4:
            return 'Пятница'
        case 5:
            return 'Суббота'
        case 6:
            return 'Воскресенье'

    }
}
export function getTimeString(time) {
    // Проверяем, является ли входная строка строкой формата "час:минуты"
    if (typeof time === 'string' && time.match(/^\d{1,2}:\d{2}$/)) {
        return time; // Возвращаем строку без изменений
    }

    // Если входная строка не в формате "час:минуты", обрабатываем ее как число времени
    const date = new Date(time);
    const hours = time[11]+time[12];
    const minutes = time[14]+time[15];
    return `${hours}:${minutes}`;
}