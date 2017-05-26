//
// Created by devil on 26.05.17.
//
// Архитектура сцены

/*
 * /Scene
 *  |- /GameObject
 *  |   |- Object
 *  |   |- Object
 *  |   |- Object
 *  |   |- Object
 *  |- /GameObject
 *  |   |- Object
 *  |   |- Object
 *  |   |- Object
 *  |   |- Object
 *
 * */

#ifndef PROJECT_CSCENE_H
#define PROJECT_CSCENE_H

#include "cFactory.h"
#include "cGameObject.h"

namespace Core {
    /*   --- Сама сцена, содержит список игровых объектов ---   */
    class cScene : public cFactoryObject, cFactory<cGameObject> {
    protected:
    public:
        cScene();
    };
}

#endif //PROJECT_CSCENE_H
