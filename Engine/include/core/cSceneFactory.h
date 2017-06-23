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
        cSceneFactory(Paranoia::Engine *engine);
        virtual ~cSceneFactory();

        // Создать новую сцену
        cScene* CreateNew(std::string name, int id = 0, bool lock = false);
        // Изменить активную сцену
        void SetActive(std::string name);
        // Получить активную сцену
        cScene* GetActive();

        void Update(int dt);
    };
}

#endif //PROJECT_CSCENEFACTORY_H
