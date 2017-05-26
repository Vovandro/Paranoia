//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;

class MyThreads: public System::cThread {
protected:
public:
    MyThreads(int id) : System::cThread(engine->threads, "names", id, true, true) {}

    void Work() override {
        System::cThread::Work();
        Lock();
        std::cout << "Threads № " << id << std::endl;
        UnLock();
        SleepThis(100);
    }
};

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    MyThreads thread(1);
    MyThreads thread2(2);
    MyThreads thread3(3);
    MyThreads thread4(4);
    MyThreads thread5(5);

    engine->threads->AddWork(&thread, true);
    engine->threads->AddWork(&thread2, true);
    engine->threads->AddWork(&thread3, true);
    engine->threads->AddWork(&thread4, true);
    engine->threads->AddWork(&thread5, true);

    engine->Start();

    delete engine;

    return 0;
}