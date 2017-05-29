//
// Created by devil on 26.05.17.
//

#include "../../include/core/cScene.h"

Core::cScene::cScene(std::string name, int id, bool lock) : Core::cFactoryObject(name, id, lock) {

}

Core::cScene::~cScene() {

}

void Core::cScene::Update(int dt) {
    for (int i = 0; i < obj.size(); i++) {
        if (obj[i])
            obj[i]->Update(dt);
    }
}
