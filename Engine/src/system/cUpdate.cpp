//
// Created by devil on 25.05.17.
//

#include "../../include/system/cUpdate.h"
#include "../../include/engine.h"

System::cUpdate::cUpdate(Paranoia::Engine *engine): cThread(engine->threads, "update", 1, true, true, 1, true) {
    this->engine = engine;
    engine->threads->AddWork(this, true);
}

System::cUpdate::~cUpdate() {
    Unlock2D();
    Unlock3D();
}

void System::cUpdate::Lock2D() {
    mutex2D.lock();
}

void System::cUpdate::Unlock2D() {
    mutex2D.unlock();
}

void System::cUpdate::Lock3D() {
    mutex3D.lock();
}

void System::cUpdate::Unlock3D() {
    mutex3D.unlock();
}

void System::cUpdate::Work() {
    cThread::Work();

    LockLocal();

    if (engine->states)
        engine->states->Update(0);

    UnLockLocal();
}
