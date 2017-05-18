//
// Created by devil on 18.05.17.
//

#include "../include/engine.h"

Paranoia::Engine::Engine(eStartType type) {
    run = false;
    this->type = type;

    window = new System::cWindow(this);
    threads = new System::cThreadFactory();
}

Paranoia::Engine::~Engine() {
    Stop();

    delete threads;
    delete window;
}


bool Paranoia::Engine::Init() {

    window->Init(3, 0, 0);

    run = true;
    return true;
}

void Paranoia::Engine::Start() {
    while (run) {
        handleEvents();
    }
}

void Paranoia::Engine::Stop() {
    run = false;
}

void Paranoia::Engine::handleEvents() {
    sf::Event event;

    if (window->GetWindow()->pollEvent(event)) {
        switch (event.type) {
            case sf::Event::Closed:
                run = false;
                break;

            case sf::Event::LostFocus:

                break;

            case sf::Event::GainedFocus:

                break;

            case sf::Event::Resized:
                break;

            default:
                break;
        }
    }
}
