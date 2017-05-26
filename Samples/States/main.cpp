//
// Created by devil on 26.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;
int i = 0;


class sMenu : public Core::cState {
public:
    sMenu() : Core::cState("menu") {};

    virtual void Start() override { engine->log->AddMessage("Start menu state", LOG_MESSAGE); };
    virtual void Update(int dt) override {};
    virtual void End() override {};
};

class sLoading : public Core::cState {
public:
    sLoading() : Core::cState("loading") {};

    virtual void Start() override { /* Add threads from loading resource */ engine->log->AddMessage("Start loading state", LOG_MESSAGE); };
    virtual void Update(int dt) override { i++; if (i >= 500) engine->states->Pop(); };
    virtual void End() override { /* End loading resource */ engine->log->AddMessage("End loading state", LOG_MESSAGE); sMenu *menuState = new sMenu(); engine->states->Push(menuState); };
};

int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    sLoading *loadingState;
    loadingState = new sLoading();

    engine->states->Push(loadingState);

    engine->Start();

    delete engine;

    return 0;
}
