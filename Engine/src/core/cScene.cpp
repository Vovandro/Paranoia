//
// Created by devil on 26.05.17.
//

#include "../../include/core/cScene.h"

Core::cScene::cScene(Paranoia::Engine *engine, std::string name, int id, bool lock) : Core::cFactoryObject(engine, name, id, lock), Core::cFactory<cGameObject>(engine) {

}

Core::cScene::~cScene() {

}

void Core::cScene::Update(int dt) {
    for (int i = 0; i < obj.size(); i++) {
        if (obj[i])
            obj[i]->Update(dt);
    }
}
