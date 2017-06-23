//
// Created by devil on 26.05.17.
//

#include "../../include/core/cGameObject.h"

Core::cGameObject::cGameObject(Paranoia::Engine *engine, std::string name, int id, bool lock) : Core::cFactoryObject(engine, name, id, lock), Core::cFactory<cGameObject>(engine) {
    objects = new Core::cFactory<cObject>(engine);
}

Core::cGameObject::~cGameObject() {
    delete objects;
}

void Core::cGameObject::Update(int dt) {
    for (int i = 0; i < obj.size(); i++) {
        if (obj[i])
            obj[i]->Update(dt);
    }

    std::vector<cObject*> *tmp = objects->GetAll();
    for (int i = 0; i < tmp->size(); i++) {
        if ((*tmp)[i])
            (*tmp)[i]->Update(dt);
    }
}