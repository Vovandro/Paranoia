//
// Created by devil on 25.05.17.
//

#include "../../include/core/cUpdate.h"
#include "../../include/engine.h"

Core::cUpdate::cUpdate(Paranoia::Engine *engine): System::cThread(engine->threads, "update", 1, true, true, 1, true) {
    this->engine = engine;
    engine->threads->AddWork(this, true);
}

Core::cUpdate::~cUpdate() {
    Unlock2D();
    Unlock3D();
}

void Core::cUpdate::Lock2D() {
    mutex2D.lock();
}

void Core::cUpdate::Unlock2D() {
    mutex2D.unlock();
}

void Core::cUpdate::Lock3D() {
    mutex3D.lock();
}

void Core::cUpdate::Unlock3D() {
    mutex3D.unlock();
}

void Core::cUpdate::Work() {
    System::cThread::Work();

    LockLocal();

    if (engine->states)
        engine->states->Update(0);

    UnLockLocal();
}
