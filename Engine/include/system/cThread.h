//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CTHREAD_H
#define PROJECT_CTHREAD_H

#include "../core/cFactoryObject.h"
#include <SFML/System.hpp>

namespace System {
    class cThreadFactory;

    class cThread : public Core::cFactoryObject {
    protected:
        cThreadFactory *factory;
        sf::Mutex localMutex;
        sf::Thread *thread;

        bool lockGlobal;
        bool enabled;
        bool loop;
        bool isStop;
        bool isMessage;
        unsigned int updateTime;
    public:
        cThread(cThreadFactory *factory, std::string name, int id, bool enabled = true, bool loop = false, unsigned int updateTime = 1, bool lock = false);
        virtual ~cThread();

        //Цикл потока
        void Thread();

        // Ассоциаци созданного потока
        void AssignThread(sf::Thread *thread);
        // Инициализаци потока
        void Init();
        // Запуск потока
        void Start();
        // Принудительная остановка потока
        void Destroy();
        // Изменение активности потока
        void SetEnabled(bool enabled);
        // Изменение параметра цикличности потока
        void SetLoop(bool loop);
        // Установка переменной дл остановки потока
        void SetStop();
        // Получение параметра остановлен ли поток
        bool GetStop();
        // Получение необходимости вызова сообщения из основного потока
        bool GetMessage();

        //Блокировка общих ресурсов необходимая для безопасной записи
        void Lock();
        //Разблокировка общих ресурсов
        void UnLock();

        //Блокировка локальных ресурсов класса для безопасной записи
        void LockLocal();
        //Разблокировка локальных ресурсов
        void UnLockLocal();

        // Заморозка потока
        void SleepThis(unsigned int time);


        //Действие выполняемое в потоке
        virtual void Work();

        //Действие выполняемое при закрытие потока
        virtual void EndWork();

        //Вызываеться при необходимости в основном потоке
        virtual void Message();
    };
}

#endif //PROJECT_CTHREAD_H
