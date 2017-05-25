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

    public:
        cUpdate(Paranoia::Engine *engine);

    };
}

#endif //PROJECT_CUPDATE_H
