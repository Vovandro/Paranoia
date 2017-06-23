//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CTHREADFACTORY_H
#define PROJECT_CTHREADFACTORY_H

#include "../core/cFactory.h"
#include "cThread.h"

namespace System {
    class cThreadFactory : public Core::cFactory<cThread> {
    protected:
    public:
        cThreadFactory(Paranoia::Engine *engine);
        virtual ~cThreadFactory();

        sf::Mutex globalMutex;

        //Добавление потока в пул
        virtual void AddObject(cThread *work, bool play);
        virtual void AddObject(cThread *work) override;
        //Продолжение приостановленного потока
        void Play(std::string name);
        //Приостановление потока на паузу
        void Pause(std::string name);
        //Остановка потока и завершение его работы
        void Stop(std::string name);
        //Принудительная остановка потока
        void Destroy(std::string name);

        // Принудительна остановка всех потоков
        void DestroyFull();

        //системное, проверка на сообщения и удаление мусора
        void Update();
    };
}

#endif //PROJECT_CTHREADFACTORY_H
