//
// Created by devil on 01.06.17.
//

#include "engine.h"

Paranoia::Engine *engine;

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init("engine.cf");

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

    Core::cConfig mk("make_config", 1);

    mk.FromString("2|?|0|||id=1==12|?|1|||name=0==Game Objects");

    engine->Start();


    delete engine;
    return 0;
}