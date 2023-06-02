## PWDM_SERVER
Серверная часть менеджера паролей. Работает про протоколу gRPC. Proto файл совместно
с клиентской частью [pwdm_client](https://github.com/BillyBones007/pwdm_client) используется
из отдельного репозитория: [pwdm_service_api](https://github.com/BillyBones007/pwdm_service_api).


#### Общая схема работы приложения
[1]: /assets/password_manager.png
![Password Manager][1]

#### Структура базы данных
[2]: /assets/data_base_arch.png
![Structure data base][2]
