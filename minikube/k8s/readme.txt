

Fixed deployment:

1. Удаляем текущий деплой:
kubectl delete -f deployment.yml

2. Правим версию докер образа в манифесте и применяем новые настройки:
kubectl apply -f deployment.yml


Rolling Update deployment:

1. Сразу правим версию докер образа в манифесте и применяем новые настройки:
kubectl apply -f deployment.yml