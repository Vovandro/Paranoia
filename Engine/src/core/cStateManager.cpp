//
// Created by devil on 25.05.17.
//

#include "../../include/core/cStateManager.h"

void Core::cStateManager::Push(cState *newState) {
    if (newState)
    {
        newState->prev = state;
        state = newState;
        state->Start();
    }
}

void Core::cStateManager::Pop() {
    if (state)
    {
        cState *last = state;
        state->End();
        state = state->prev;
        delete last;
    }
}

void Core::cStateManager::PopAll(bool isMessage) {
    while (state) {
        cState *last = state;
        if (isMessage)
            state->End();
        state = state->prev;
        delete last;
    }
}

Core::cState *Core::cStateManager::Get() {
    return state;
}

void Core::cStateManager::Update(int dt) {
    if (state)
        state->Update(dt);
}

Core::cStateManager::cStateManager() {
    state = NULL;
}


Core::cStateManager::~cStateManager() {
    PopAll(false);
}