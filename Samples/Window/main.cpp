//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    engine->Start();


    delete engine;
    return 0;
}