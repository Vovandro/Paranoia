//
// Created by devil on 25.05.17.
//

#ifndef PROJECT_CUPDATE_H
#define PROJECT_CUPDATE_H

#include "../system/cThread.h"

namespace Paranoia {
    class Engine;
}

namespace Core {
    class cUpdate : public System::cThread {
    protected:
        Paranoia::Engine *engine;
        sf::Mutex mutex2D;
        sf::Mutex mutex3D;

    public:
        cUpdate(Paranoia::Engine *engine);
        ~cUpdate();

        // Блокировка 2Д данных для потоков
        void Lock2D();
        // Разблокировка 2Д данных дл потоков
        void Unlock2D();
        // Блокировка 3Д данных для потоков
        void Lock3D();
        // Разблокировка 3Д данных для потоков
        void Unlock3D();

        /* поток обновления */
        virtual void Work() override;
    };
}

#endif //PROJECT_CUPDATE_H
