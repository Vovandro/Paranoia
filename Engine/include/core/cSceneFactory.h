//
// Created by devil on 29.05.17.
//

#ifndef PROJECT_CSCENEFACTORY_H
#define PROJECT_CSCENEFACTORY_H

#include "cFactory.h"
#include "cScene.h"

namespace Core {
    class cSceneFactory : public cFactory<cScene> {
    protected:
        cScene *activeScene;

    public:
        cSceneFactory();
        virtual ~cSceneFactory();

        cScene* CreateNew(std::string name, int id = 0, bool lock = false);
        void SetActive(std::string name);
        cScene* GetActive();

        void Update(int dt);
    };
}

#endif //PROJECT_CSCENEFACTORY_H
