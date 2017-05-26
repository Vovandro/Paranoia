//
// Created by devil on 25.05.17.
//

#include "../../include/render/cRender.h"
#include "../../include/engine.h"

Render::cRender::cRender(Paranoia::Engine *engine) {
    this->engine = engine;
}

Render::cRender::~cRender() {
}

bool Render::cRender::Init() {
    return false;
}

void Render::cRender::Update() {
    engine->window->GetWindow()->setActive(true);

    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

    glColor3ub(255,0,0);
    glBegin(GL_QUADS);
    glVertex2f(50,400);
    glVertex2f(50,50);
    glVertex2f(200,50);
    glVertex2f(200,200);
    glEnd();
}

void Render::cRender::Resize(int w, int h) {
    glViewport(0, 0, w, h);
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    glOrtho(0, w, 0, h, -1.0, 1.0);
    glMatrixMode(GL_MODELVIEW);
    glLoadIdentity();
}