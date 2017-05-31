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
    cScene *ret = new cScene(name, id, lock);

    AddObject(ret);

    return ret;
}

void Core::cSceneFactory::SetActive(std::string name) {
    cScene *tmp = FindObject(name);

    if (tmp) {
        activeScene = tmp;
    }
}

Core::cScene *Core::cSceneFactory::GetActive() {
    return activeScene;
}

void Core::cSceneFactory::Update(int dt) {
    for (int i = 0; i < obj.size(); i++) {
        if (obj[i]) {
            obj[i]->Update(dt);
        }
    }
}



