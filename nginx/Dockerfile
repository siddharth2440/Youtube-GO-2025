FROM nginx:1.21.6
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# this commnad ensures that nginx will run on the foreground instead of background
CMD [ "nginx","-g","daemon off;" ]