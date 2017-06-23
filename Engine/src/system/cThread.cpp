//
// Created by devil on 18.05.17.
//

#include "../../include/system/cThread.h"
#include "../../include/system/cThreadFactory.h"
#include "../../include/engine.h"

System::cThread::cThread(Paranoia::Engine *engine, std::string name, int id, bool enabled, bool loop, unsigned int updateTime, bool lock) : Core::cFactoryObject(engine, name, id, lock) {
    this->enabled = enabled;
    this->loop = loop;
    this->updateTime = updateTime;

    isMessage = false;
    isStop = false;
    lockGlobal = false;
}

System::cThread::~cThread() {
    if (lockGlobal)
        UnLock();

    UnLockLocal();
}

void System::cThread::Thread() {
    do {
        while (!enabled) {
            SleepThis(updateTime);

            if (isStop) {
                EndWork();
                return;
            }
        }

        Work();
    } while (loop);

    EndWork();
}

void System::cThread::AssignThread(sf::Thread *thread) {
    this->thread = thread;
}

void System::cThread::Init() {
    thread = new sf::Thread(&cThread::Thread, this);
}

void System::cThread::Start() {
    enabled = true;
    if (thread)
        thread->launch();
}

void System::cThread::Destroy() {
    LockLocal();

    if (thread)
        thread->terminate();

    if (lockGlobal)
        UnLock();

    UnLockLocal();
}

void System::cThread::SetEnabled(bool enabled) {
    LockLocal();
    this->enabled = enabled;
    UnLockLocal();
}

void System::cThread::SetLoop(bool loop) {
    LockLocal();
    this->loop = loop;
    UnLockLocal();
}

void System::cThread::SetStop() {
    LockLocal();
    isStop = true;
    UnLockLocal();
}

bool System::cThread::GetStop() {
    return isStop;
}

bool System::cThread::GetMessage() {
    return isMessage;
}

void System::cThread::Lock() {
    engine->threads->globalMutex.lock();
    lockGlobal = true;
}

void System::cThread::UnLock() {
    engine->threads->globalMutex.unlock();
    lockGlobal = false;
}

void System::cThread::LockLocal() {
    localMutex.lock();
}

void System::cThread::UnLockLocal() {
    localMutex.unlock();
}

void System::cThread::SleepThis(unsigned int time) {
    sf::sleep(sf::milliseconds(time));
}

void System::cThread::Work() {
    SleepThis(1);
}

void System::cThread::EndWork() {
    SetStop();
}

void System::cThread::Message() {
}

void System::cThread::Register() {
    if (engine) {
        engine->threads->AddObject(this, false);
    }
}
