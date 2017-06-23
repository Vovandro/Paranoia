//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;

class MyThreads: public System::cThread {
protected:
public:
    MyThreads(int id) : System::cThread(engine, "names", id, true, true) {}

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

    engine->Init("engine.cf");

    MyThreads thread(1);
    MyThreads thread2(2);
    MyThreads thread3(3);
    MyThreads thread4(4);
    MyThreads thread5(5);

    engine->threads->AddObject(&thread, true);
    engine->threads->AddObject(&thread2, true);
    engine->threads->AddObject(&thread3, true);
    engine->threads->AddObject(&thread4, true);
    engine->threads->AddObject(&thread5, true);

    engine->Start();

    delete engine;

    return 0;
}