//
// Created by devil on 29.05.17.
//

#include "../../include/core/cSceneFactory.h"

Core::cSceneFactory::cSceneFactory() {
    activeScene = NULL;
}

Core::cSceneFactory::~cSceneFactory() {

}

Core::cScene *Core::cSceneFactory::CreateNew(std::string name, int id, bool lock) {
    return NULL;
}

void Core::cSceneFactory::SetActive(std::string name) {

}

Core::cScene *Core::cSceneFactory::GetActive() {
    return NULL;
}

void Core::cSceneFactory::Update(int dt) {
    for (int i = 0; i < obj.size(); i++) {
        if (obj[i]) {
            obj[i]->Update(dt);
        }
    }
}



