export function getJwtToken() {
    let token = localStorage.getItem('token')
    return token

}
export function getTokenPayload(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}
export function getTokenValues(token){
    const decodedToken = getTokenPayload(token);
    const email = decodedToken.email;
    const name = decodedToken.name;

    return {
        email: email,
        name: name
    };
}

export function getUserRoleFromToken(token){
    // Получаем токен из localStorage
    // Если токен не существует, возвращаем null
    if (!token) {
        return null;
    }
    try {
        // Декодируем токен из Base64
        let tokenPayload = JSON.parse(atob(token.split('.')[1]));
        // Извлекаем роль пользователя из токена
        let role = tokenPayload.role;
        return role;
    } catch (error) {
        console.error('Ошибка при декодировании токена:', error);
        return null;
    }
}
export function removeToken(){
    localStorage.removeItem("token");
}