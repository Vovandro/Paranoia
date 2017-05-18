//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;

class MyThreads: public System::cThread {
protected:
public:
    MyThreads() : System::cThread(engine->threads, "names", 10, true, true) {}

    void Work() override {
        System::cThread::Work();
        LockLocal();
        std::cout << "Threads" << std::endl;
        UnLockLocal();
    }
};

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    MyThreads thread;
    MyThreads thread2;
    MyThreads thread3;

    engine->threads->AddWork(&thread, true);

    engine->threads->AddWork(&thread2, true);

    engine->threads->AddWork(&thread3, true);

    engine->Start();

    return 0;
}