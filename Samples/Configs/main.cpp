//
// Created by devil on 01.06.17.
//

#include "engine.h"

Paranoia::Engine *engine;

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    Core::cConfig conf("test", 1);

    Core::cConfigItemInt iId;

    iId.name = "id";
    iId.data = 12;

    Core::cConfigItemString iName;

    iName.name = "name";
    iName.data = "Game Objects";


    conf.Add(&iId);
    conf.Add(&iName);

    std::cout << conf.ToString();

    engine->Start();


    delete engine;
    return 0;
}