import * as tokenFunctions from './pakages/token.js';
let cur_page = 0;
const limit = 10;
const activitiesList = document.getElementById("activities-list");
const profileElement = document.getElementById("profile-details");

document.addEventListener("DOMContentLoaded", function(){
    let token = tokenFunctions.getJwtToken()
    alert(token)
})