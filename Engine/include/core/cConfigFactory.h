//
// Created by devil on 17.06.17.
//

#ifndef PROJECT_CCONFIGFACTORY_H
#define PROJECT_CCONFIGFACTORY_H

#include "cConfig.h"
#include "cFactory.h"


namespace Core {
    class cConfigFactory : public cFactory<cConfig> {
    protected:
    public:

        cConfigFactory(Paranoia::Engine *engine);
        virtual ~cConfigFactory();

        virtual void AddObject(cConfig* newObj) override;
        virtual cConfig* AddObject(std::string cfName, int id = 0, bool lock = false);
        //virtual cConfig* AddObject(std::string afName, std::string cfName, int id = 0, bool lock = false);
        //virtual cConfig* CreateObject(std::string config, std::string name, int id = 0, bool lock = false);
    };
}

#endif //PROJECT_CCONFIGFACTORY_H
