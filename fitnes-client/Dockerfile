FROM nginx:latest

# Копируем все файлы приложения в папку nginx
COPY . /usr/share/nginx/html

# Создаем директории для статических ресурсов
RUN mkdir -p /usr/share/nginx/html/scripts
RUN mkdir -p /usr/share/nginx/html/styles
RUN mkdir -p /usr/share/nginx/html/src

# Копируем статические файлы в соответствующие папки
COPY scripts/* /usr/share/nginx/html/scripts
COPY styles/* /usr/share/nginx/html/styles
COPY src/* /usr/share/nginx/html/src

# Определяем порт
EXPOSE 63342

# Создаем файл конфигурации Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]