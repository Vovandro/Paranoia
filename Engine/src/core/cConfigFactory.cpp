//
// Created by devil on 17.06.17.
//

#include "../../include/core/cConfigFactory.h"
#include "../../include/engine.h"

Core::cConfigFactory::cConfigFactory(Paranoia::Engine *engine) : Core::cFactory<cConfig>(engine) {

}

Core::cConfigFactory::~cConfigFactory() {

}

Core::cConfig* Core::cConfigFactory::AddObject(std::string cfName, int id, bool lock) {
    cConfig *newItem = NULL;

    if (engine) {
        if (engine->files->Open(cfName, OPEN_READ)) {
            newItem = new cConfig(engine, cfName, (id==0)?GetNewID():id, lock);
            newItem->FromString(engine->files->Read(cfName)->toStr());

            AddObject(newItem);
        }
    }

    return newItem;
}

void Core::cConfigFactory::AddObject(Core::cConfig *newObj) {
    cFactory::AddObject(newObj);
}
