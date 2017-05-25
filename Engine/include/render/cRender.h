//
// Created by devil on 25.05.17.
//

#ifndef PROJECT_CRENDER_H
#define PROJECT_CRENDER_H

#include <GL/gl.h>

namespace Paranoia {
    class Engine;
}

namespace Render {
    class cRender {
    protected:
        Paranoia::Engine *engine;

    public:
        cRender(Paranoia::Engine *engine);
        ~cRender();

        bool Init();
        void Update(float dt);
        void Resize(int w, int h);
    };
}

#endif //PROJECT_CRENDER_H
