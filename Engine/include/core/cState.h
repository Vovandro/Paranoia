//
// Created by devil on 25.05.17.
//

#ifndef PROJECT_CSTATE_H
#define PROJECT_CSTATE_H

#include <string>

namespace Core {
    class cState {
    public:
        cState *prev;
        std::string id;

        cState(std::string id);
        virtual ~cState();

        virtual void Start() = 0;
        virtual void Update(int dt) = 0;
        virtual void End() = 0;
    };
}

#endif //PROJECT_CSTATE_H
