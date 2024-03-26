FROM python:3.10.12

WORKDIR /admin

COPY requirements.txt ./

RUN pip install -r requirements.txt

COPY . ./

EXPOSE 9000
