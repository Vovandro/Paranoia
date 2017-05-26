//
// Created by devil on 25.05.17.
//

#ifndef PROJECT_CUPDATE_H
#define PROJECT_CUPDATE_H

#include "cThread.h"

namespace Paranoia {
    class Engine;
}

namespace System {
    class cUpdate : public cThread {
    protected:
        Paranoia::Engine *engine;
        sf::Mutex mutex2D;
        sf::Mutex mutex3D;

    public:
        cUpdate(Paranoia::Engine *engine);
        ~cUpdate();

        void Lock2D();
        void Unlock2D();
        void Lock3D();
        void Unlock3D();

        virtual void Work() override;
    };
}

#endif //PROJECT_CUPDATE_H
