/****************************************************************************
** Meta object code from reading C++ file 'moc.cpp'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.13.0)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include <memory>
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'moc.cpp' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.13.0. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
QT_WARNING_PUSH
QT_WARNING_DISABLE_DEPRECATED
struct qt_meta_stringdata_ApproveListingCtx687eda_t {
    QByteArrayData data[23];
    char stringdata0[242];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_ApproveListingCtx687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_ApproveListingCtx687eda_t qt_meta_stringdata_ApproveListingCtx687eda = {
    {
QT_MOC_LITERAL(0, 0, 23), // "ApproveListingCtx687eda"
QT_MOC_LITERAL(1, 24, 4), // "back"
QT_MOC_LITERAL(2, 29, 0), // ""
QT_MOC_LITERAL(3, 30, 7), // "approve"
QT_MOC_LITERAL(4, 38, 6), // "reject"
QT_MOC_LITERAL(5, 45, 19), // "onCheckStateChanged"
QT_MOC_LITERAL(6, 65, 1), // "i"
QT_MOC_LITERAL(7, 67, 7), // "checked"
QT_MOC_LITERAL(8, 75, 13), // "triggerUpdate"
QT_MOC_LITERAL(9, 89, 13), // "remoteChanged"
QT_MOC_LITERAL(10, 103, 6), // "remote"
QT_MOC_LITERAL(11, 110, 16), // "transportChanged"
QT_MOC_LITERAL(12, 127, 9), // "transport"
QT_MOC_LITERAL(13, 137, 15), // "endpointChanged"
QT_MOC_LITERAL(14, 153, 8), // "endpoint"
QT_MOC_LITERAL(15, 162, 11), // "fromChanged"
QT_MOC_LITERAL(16, 174, 4), // "from"
QT_MOC_LITERAL(17, 179, 14), // "messageChanged"
QT_MOC_LITERAL(18, 194, 7), // "message"
QT_MOC_LITERAL(19, 202, 14), // "rawDataChanged"
QT_MOC_LITERAL(20, 217, 7), // "rawData"
QT_MOC_LITERAL(21, 225, 11), // "hashChanged"
QT_MOC_LITERAL(22, 237, 4) // "hash"

    },
    "ApproveListingCtx687eda\0back\0\0approve\0"
    "reject\0onCheckStateChanged\0i\0checked\0"
    "triggerUpdate\0remoteChanged\0remote\0"
    "transportChanged\0transport\0endpointChanged\0"
    "endpoint\0fromChanged\0from\0messageChanged\0"
    "message\0rawDataChanged\0rawData\0"
    "hashChanged\0hash"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_ApproveListingCtx687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
      12,   14, // methods
       7,  104, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
      12,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   74,    2, 0x06 /* Public */,
       3,    0,   75,    2, 0x06 /* Public */,
       4,    0,   76,    2, 0x06 /* Public */,
       5,    2,   77,    2, 0x06 /* Public */,
       8,    0,   82,    2, 0x06 /* Public */,
       9,    1,   83,    2, 0x06 /* Public */,
      11,    1,   86,    2, 0x06 /* Public */,
      13,    1,   89,    2, 0x06 /* Public */,
      15,    1,   92,    2, 0x06 /* Public */,
      17,    1,   95,    2, 0x06 /* Public */,
      19,    1,   98,    2, 0x06 /* Public */,
      21,    1,  101,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void, QMetaType::Int, QMetaType::Bool,    6,    7,
    QMetaType::Void,
    QMetaType::Void, QMetaType::QString,   10,
    QMetaType::Void, QMetaType::QString,   12,
    QMetaType::Void, QMetaType::QString,   14,
    QMetaType::Void, QMetaType::QString,   16,
    QMetaType::Void, QMetaType::QString,   18,
    QMetaType::Void, QMetaType::QString,   20,
    QMetaType::Void, QMetaType::QString,   22,

 // properties: name, type, flags
      10, QMetaType::QString, 0x00495103,
      12, QMetaType::QString, 0x00495103,
      14, QMetaType::QString, 0x00495103,
      16, QMetaType::QString, 0x00495103,
      18, QMetaType::QString, 0x00495103,
      20, QMetaType::QString, 0x00495103,
      22, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       5,
       6,
       7,
       8,
       9,
      10,
      11,

       0        // eod
};

void ApproveListingCtx687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<ApproveListingCtx687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->back(); break;
        case 1: _t->approve(); break;
        case 2: _t->reject(); break;
        case 3: _t->onCheckStateChanged((*reinterpret_cast< qint32(*)>(_a[1])),(*reinterpret_cast< bool(*)>(_a[2]))); break;
        case 4: _t->triggerUpdate(); break;
        case 5: _t->remoteChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 6: _t->transportChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 7: _t->endpointChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 8: _t->fromChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 9: _t->messageChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 10: _t->rawDataChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 11: _t->hashChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (ApproveListingCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::back)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::approve)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::reject)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(qint32 , bool );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::onCheckStateChanged)) {
                *result = 3;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::triggerUpdate)) {
                *result = 4;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::remoteChanged)) {
                *result = 5;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::transportChanged)) {
                *result = 6;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::endpointChanged)) {
                *result = 7;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::fromChanged)) {
                *result = 8;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::messageChanged)) {
                *result = 9;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::rawDataChanged)) {
                *result = 10;
                return;
            }
        }
        {
            using _t = void (ApproveListingCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveListingCtx687eda::hashChanged)) {
                *result = 11;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<ApproveListingCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->remote(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->transport(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->endpoint(); break;
        case 3: *reinterpret_cast< QString*>(_v) = _t->from(); break;
        case 4: *reinterpret_cast< QString*>(_v) = _t->message(); break;
        case 5: *reinterpret_cast< QString*>(_v) = _t->rawData(); break;
        case 6: *reinterpret_cast< QString*>(_v) = _t->hash(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<ApproveListingCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setRemote(*reinterpret_cast< QString*>(_v)); break;
        case 1: _t->setTransport(*reinterpret_cast< QString*>(_v)); break;
        case 2: _t->setEndpoint(*reinterpret_cast< QString*>(_v)); break;
        case 3: _t->setFrom(*reinterpret_cast< QString*>(_v)); break;
        case 4: _t->setMessage(*reinterpret_cast< QString*>(_v)); break;
        case 5: _t->setRawData(*reinterpret_cast< QString*>(_v)); break;
        case 6: _t->setHash(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject ApproveListingCtx687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_ApproveListingCtx687eda.data,
    qt_meta_data_ApproveListingCtx687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *ApproveListingCtx687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *ApproveListingCtx687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_ApproveListingCtx687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int ApproveListingCtx687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 12)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 12;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 12)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 12;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 7;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 7;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 7;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 7;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 7;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 7;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void ApproveListingCtx687eda::back()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void ApproveListingCtx687eda::approve()
{
    QMetaObject::activate(this, &staticMetaObject, 1, nullptr);
}

// SIGNAL 2
void ApproveListingCtx687eda::reject()
{
    QMetaObject::activate(this, &staticMetaObject, 2, nullptr);
}

// SIGNAL 3
void ApproveListingCtx687eda::onCheckStateChanged(qint32 _t1, bool _t2)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))), const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t2))) };
    QMetaObject::activate(this, &staticMetaObject, 3, _a);
}

// SIGNAL 4
void ApproveListingCtx687eda::triggerUpdate()
{
    QMetaObject::activate(this, &staticMetaObject, 4, nullptr);
}

// SIGNAL 5
void ApproveListingCtx687eda::remoteChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 5, _a);
}

// SIGNAL 6
void ApproveListingCtx687eda::transportChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 6, _a);
}

// SIGNAL 7
void ApproveListingCtx687eda::endpointChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 7, _a);
}

// SIGNAL 8
void ApproveListingCtx687eda::fromChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 8, _a);
}

// SIGNAL 9
void ApproveListingCtx687eda::messageChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 9, _a);
}

// SIGNAL 10
void ApproveListingCtx687eda::rawDataChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 10, _a);
}

// SIGNAL 11
void ApproveListingCtx687eda::hashChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 11, _a);
}
struct qt_meta_stringdata_ApproveNewAccountCtx687eda_t {
    QByteArrayData data[17];
    char stringdata0[217];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_ApproveNewAccountCtx687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_ApproveNewAccountCtx687eda_t qt_meta_stringdata_ApproveNewAccountCtx687eda = {
    {
QT_MOC_LITERAL(0, 0, 26), // "ApproveNewAccountCtx687eda"
QT_MOC_LITERAL(1, 27, 7), // "clicked"
QT_MOC_LITERAL(2, 35, 0), // ""
QT_MOC_LITERAL(3, 36, 1), // "b"
QT_MOC_LITERAL(4, 38, 4), // "back"
QT_MOC_LITERAL(5, 43, 14), // "passwordEdited"
QT_MOC_LITERAL(6, 58, 21), // "confirmPasswordEdited"
QT_MOC_LITERAL(7, 80, 13), // "remoteChanged"
QT_MOC_LITERAL(8, 94, 6), // "remote"
QT_MOC_LITERAL(9, 101, 16), // "transportChanged"
QT_MOC_LITERAL(10, 118, 9), // "transport"
QT_MOC_LITERAL(11, 128, 15), // "endpointChanged"
QT_MOC_LITERAL(12, 144, 8), // "endpoint"
QT_MOC_LITERAL(13, 153, 15), // "passwordChanged"
QT_MOC_LITERAL(14, 169, 8), // "password"
QT_MOC_LITERAL(15, 178, 22), // "confirmPasswordChanged"
QT_MOC_LITERAL(16, 201, 15) // "confirmPassword"

    },
    "ApproveNewAccountCtx687eda\0clicked\0\0"
    "b\0back\0passwordEdited\0confirmPasswordEdited\0"
    "remoteChanged\0remote\0transportChanged\0"
    "transport\0endpointChanged\0endpoint\0"
    "passwordChanged\0password\0"
    "confirmPasswordChanged\0confirmPassword"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_ApproveNewAccountCtx687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       9,   14, // methods
       5,   84, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       9,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   59,    2, 0x06 /* Public */,
       4,    0,   62,    2, 0x06 /* Public */,
       5,    1,   63,    2, 0x06 /* Public */,
       6,    1,   66,    2, 0x06 /* Public */,
       7,    1,   69,    2, 0x06 /* Public */,
       9,    1,   72,    2, 0x06 /* Public */,
      11,    1,   75,    2, 0x06 /* Public */,
      13,    1,   78,    2, 0x06 /* Public */,
      15,    1,   81,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, QMetaType::Int,    3,
    QMetaType::Void,
    QMetaType::Void, QMetaType::QString,    3,
    QMetaType::Void, QMetaType::QString,    3,
    QMetaType::Void, QMetaType::QString,    8,
    QMetaType::Void, QMetaType::QString,   10,
    QMetaType::Void, QMetaType::QString,   12,
    QMetaType::Void, QMetaType::QString,   14,
    QMetaType::Void, QMetaType::QString,   16,

 // properties: name, type, flags
       8, QMetaType::QString, 0x00495103,
      10, QMetaType::QString, 0x00495103,
      12, QMetaType::QString, 0x00495103,
      14, QMetaType::QString, 0x00495103,
      16, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       4,
       5,
       6,
       7,
       8,

       0        // eod
};

void ApproveNewAccountCtx687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<ApproveNewAccountCtx687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->clicked((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 1: _t->back(); break;
        case 2: _t->passwordEdited((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 3: _t->confirmPasswordEdited((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 4: _t->remoteChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 5: _t->transportChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 6: _t->endpointChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 7: _t->passwordChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 8: _t->confirmPasswordChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::clicked)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::back)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::passwordEdited)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::confirmPasswordEdited)) {
                *result = 3;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::remoteChanged)) {
                *result = 4;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::transportChanged)) {
                *result = 5;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::endpointChanged)) {
                *result = 6;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::passwordChanged)) {
                *result = 7;
                return;
            }
        }
        {
            using _t = void (ApproveNewAccountCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveNewAccountCtx687eda::confirmPasswordChanged)) {
                *result = 8;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<ApproveNewAccountCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->remote(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->transport(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->endpoint(); break;
        case 3: *reinterpret_cast< QString*>(_v) = _t->password(); break;
        case 4: *reinterpret_cast< QString*>(_v) = _t->confirmPassword(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<ApproveNewAccountCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setRemote(*reinterpret_cast< QString*>(_v)); break;
        case 1: _t->setTransport(*reinterpret_cast< QString*>(_v)); break;
        case 2: _t->setEndpoint(*reinterpret_cast< QString*>(_v)); break;
        case 3: _t->setPassword(*reinterpret_cast< QString*>(_v)); break;
        case 4: _t->setConfirmPassword(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject ApproveNewAccountCtx687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_ApproveNewAccountCtx687eda.data,
    qt_meta_data_ApproveNewAccountCtx687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *ApproveNewAccountCtx687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *ApproveNewAccountCtx687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_ApproveNewAccountCtx687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int ApproveNewAccountCtx687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 9)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 9;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 9)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 9;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 5;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 5;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 5;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 5;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 5;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 5;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void ApproveNewAccountCtx687eda::clicked(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}

// SIGNAL 1
void ApproveNewAccountCtx687eda::back()
{
    QMetaObject::activate(this, &staticMetaObject, 1, nullptr);
}

// SIGNAL 2
void ApproveNewAccountCtx687eda::passwordEdited(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}

// SIGNAL 3
void ApproveNewAccountCtx687eda::confirmPasswordEdited(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 3, _a);
}

// SIGNAL 4
void ApproveNewAccountCtx687eda::remoteChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 4, _a);
}

// SIGNAL 5
void ApproveNewAccountCtx687eda::transportChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 5, _a);
}

// SIGNAL 6
void ApproveNewAccountCtx687eda::endpointChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 6, _a);
}

// SIGNAL 7
void ApproveNewAccountCtx687eda::passwordChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 7, _a);
}

// SIGNAL 8
void ApproveNewAccountCtx687eda::confirmPasswordChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 8, _a);
}
struct qt_meta_stringdata_TxListCtx687eda_t {
    QByteArrayData data[8];
    char stringdata0[95];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_TxListCtx687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_TxListCtx687eda_t qt_meta_stringdata_TxListCtx687eda = {
    {
QT_MOC_LITERAL(0, 0, 15), // "TxListCtx687eda"
QT_MOC_LITERAL(1, 16, 7), // "clicked"
QT_MOC_LITERAL(2, 24, 0), // ""
QT_MOC_LITERAL(3, 25, 1), // "b"
QT_MOC_LITERAL(4, 27, 21), // "shortenAddressChanged"
QT_MOC_LITERAL(5, 49, 14), // "shortenAddress"
QT_MOC_LITERAL(6, 64, 18), // "selectedSrcChanged"
QT_MOC_LITERAL(7, 83, 11) // "selectedSrc"

    },
    "TxListCtx687eda\0clicked\0\0b\0"
    "shortenAddressChanged\0shortenAddress\0"
    "selectedSrcChanged\0selectedSrc"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_TxListCtx687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       3,   14, // methods
       2,   38, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       3,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   29,    2, 0x06 /* Public */,
       4,    1,   32,    2, 0x06 /* Public */,
       6,    1,   35,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, QMetaType::Int,    3,
    QMetaType::Void, QMetaType::QString,    5,
    QMetaType::Void, QMetaType::QString,    7,

 // properties: name, type, flags
       5, QMetaType::QString, 0x00495103,
       7, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       1,
       2,

       0        // eod
};

void TxListCtx687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<TxListCtx687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->clicked((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 1: _t->shortenAddressChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 2: _t->selectedSrcChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (TxListCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListCtx687eda::clicked)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (TxListCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListCtx687eda::shortenAddressChanged)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (TxListCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListCtx687eda::selectedSrcChanged)) {
                *result = 2;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<TxListCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->shortenAddress(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->selectedSrc(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<TxListCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setShortenAddress(*reinterpret_cast< QString*>(_v)); break;
        case 1: _t->setSelectedSrc(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject TxListCtx687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_TxListCtx687eda.data,
    qt_meta_data_TxListCtx687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *TxListCtx687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *TxListCtx687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_TxListCtx687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int TxListCtx687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 3)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 3;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 3)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 3;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 2;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 2;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 2;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 2;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 2;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 2;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void TxListCtx687eda::clicked(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}

// SIGNAL 1
void TxListCtx687eda::shortenAddressChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 1, _a);
}

// SIGNAL 2
void TxListCtx687eda::selectedSrcChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}
struct qt_meta_stringdata_TxListModel687eda_t {
    QByteArrayData data[10];
    char stringdata0[73];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_TxListModel687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_TxListModel687eda_t qt_meta_stringdata_TxListModel687eda = {
    {
QT_MOC_LITERAL(0, 0, 17), // "TxListModel687eda"
QT_MOC_LITERAL(1, 18, 5), // "clear"
QT_MOC_LITERAL(2, 24, 0), // ""
QT_MOC_LITERAL(3, 25, 3), // "add"
QT_MOC_LITERAL(4, 29, 8), // "quintptr"
QT_MOC_LITERAL(5, 38, 2), // "tx"
QT_MOC_LITERAL(6, 41, 6), // "remove"
QT_MOC_LITERAL(7, 48, 1), // "i"
QT_MOC_LITERAL(8, 50, 14), // "isEmptyChanged"
QT_MOC_LITERAL(9, 65, 7) // "isEmpty"

    },
    "TxListModel687eda\0clear\0\0add\0quintptr\0"
    "tx\0remove\0i\0isEmptyChanged\0isEmpty"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_TxListModel687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       4,   14, // methods
       1,   44, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       4,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   34,    2, 0x06 /* Public */,
       3,    1,   35,    2, 0x06 /* Public */,
       6,    1,   38,    2, 0x06 /* Public */,
       8,    1,   41,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void, 0x80000000 | 4,    5,
    QMetaType::Void, QMetaType::Int,    7,
    QMetaType::Void, QMetaType::Bool,    9,

 // properties: name, type, flags
       9, QMetaType::Bool, 0x00495103,

 // properties: notify_signal_id
       3,

       0        // eod
};

void TxListModel687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<TxListModel687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->clear(); break;
        case 1: _t->add((*reinterpret_cast< quintptr(*)>(_a[1]))); break;
        case 2: _t->remove((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 3: _t->isEmptyChanged((*reinterpret_cast< bool(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (TxListModel687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListModel687eda::clear)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (TxListModel687eda::*)(quintptr );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListModel687eda::add)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (TxListModel687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListModel687eda::remove)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (TxListModel687eda::*)(bool );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListModel687eda::isEmptyChanged)) {
                *result = 3;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<TxListModel687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< bool*>(_v) = _t->isEmpty(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<TxListModel687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setIsEmpty(*reinterpret_cast< bool*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject TxListModel687eda::staticMetaObject = { {
    &QAbstractListModel::staticMetaObject,
    qt_meta_stringdata_TxListModel687eda.data,
    qt_meta_data_TxListModel687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *TxListModel687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *TxListModel687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_TxListModel687eda.stringdata0))
        return static_cast<void*>(this);
    return QAbstractListModel::qt_metacast(_clname);
}

int TxListModel687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QAbstractListModel::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 4)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 4;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 4)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 4;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 1;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void TxListModel687eda::clear()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void TxListModel687eda::add(quintptr _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 1, _a);
}

// SIGNAL 2
void TxListModel687eda::remove(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}

// SIGNAL 3
void TxListModel687eda::isEmptyChanged(bool _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 3, _a);
}
struct qt_meta_stringdata_ApproveSignDataCtx687eda_t {
    QByteArrayData data[27];
    char stringdata0[276];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_ApproveSignDataCtx687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_ApproveSignDataCtx687eda_t qt_meta_stringdata_ApproveSignDataCtx687eda = {
    {
QT_MOC_LITERAL(0, 0, 24), // "ApproveSignDataCtx687eda"
QT_MOC_LITERAL(1, 25, 7), // "clicked"
QT_MOC_LITERAL(2, 33, 0), // ""
QT_MOC_LITERAL(3, 34, 1), // "b"
QT_MOC_LITERAL(4, 36, 6), // "onBack"
QT_MOC_LITERAL(5, 43, 9), // "onApprove"
QT_MOC_LITERAL(6, 53, 8), // "onReject"
QT_MOC_LITERAL(7, 62, 6), // "edited"
QT_MOC_LITERAL(8, 69, 5), // "value"
QT_MOC_LITERAL(9, 75, 13), // "remoteChanged"
QT_MOC_LITERAL(10, 89, 6), // "remote"
QT_MOC_LITERAL(11, 96, 16), // "transportChanged"
QT_MOC_LITERAL(12, 113, 9), // "transport"
QT_MOC_LITERAL(13, 123, 15), // "endpointChanged"
QT_MOC_LITERAL(14, 139, 8), // "endpoint"
QT_MOC_LITERAL(15, 148, 11), // "fromChanged"
QT_MOC_LITERAL(16, 160, 4), // "from"
QT_MOC_LITERAL(17, 165, 14), // "messageChanged"
QT_MOC_LITERAL(18, 180, 7), // "message"
QT_MOC_LITERAL(19, 188, 14), // "rawDataChanged"
QT_MOC_LITERAL(20, 203, 7), // "rawData"
QT_MOC_LITERAL(21, 211, 11), // "hashChanged"
QT_MOC_LITERAL(22, 223, 4), // "hash"
QT_MOC_LITERAL(23, 228, 15), // "passwordChanged"
QT_MOC_LITERAL(24, 244, 8), // "password"
QT_MOC_LITERAL(25, 253, 14), // "fromSrcChanged"
QT_MOC_LITERAL(26, 268, 7) // "fromSrc"

    },
    "ApproveSignDataCtx687eda\0clicked\0\0b\0"
    "onBack\0onApprove\0onReject\0edited\0value\0"
    "remoteChanged\0remote\0transportChanged\0"
    "transport\0endpointChanged\0endpoint\0"
    "fromChanged\0from\0messageChanged\0message\0"
    "rawDataChanged\0rawData\0hashChanged\0"
    "hash\0passwordChanged\0password\0"
    "fromSrcChanged\0fromSrc"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_ApproveSignDataCtx687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
      14,   14, // methods
       9,  122, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
      14,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   84,    2, 0x06 /* Public */,
       4,    0,   87,    2, 0x06 /* Public */,
       5,    0,   88,    2, 0x06 /* Public */,
       6,    0,   89,    2, 0x06 /* Public */,
       7,    2,   90,    2, 0x06 /* Public */,
       9,    1,   95,    2, 0x06 /* Public */,
      11,    1,   98,    2, 0x06 /* Public */,
      13,    1,  101,    2, 0x06 /* Public */,
      15,    1,  104,    2, 0x06 /* Public */,
      17,    1,  107,    2, 0x06 /* Public */,
      19,    1,  110,    2, 0x06 /* Public */,
      21,    1,  113,    2, 0x06 /* Public */,
      23,    1,  116,    2, 0x06 /* Public */,
      25,    1,  119,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, QMetaType::Int,    3,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void, QMetaType::QString, QMetaType::QString,    3,    8,
    QMetaType::Void, QMetaType::QString,   10,
    QMetaType::Void, QMetaType::QString,   12,
    QMetaType::Void, QMetaType::QString,   14,
    QMetaType::Void, QMetaType::QString,   16,
    QMetaType::Void, QMetaType::QString,   18,
    QMetaType::Void, QMetaType::QString,   20,
    QMetaType::Void, QMetaType::QString,   22,
    QMetaType::Void, QMetaType::QString,   24,
    QMetaType::Void, QMetaType::QString,   26,

 // properties: name, type, flags
      10, QMetaType::QString, 0x00495103,
      12, QMetaType::QString, 0x00495103,
      14, QMetaType::QString, 0x00495103,
      16, QMetaType::QString, 0x00495103,
      18, QMetaType::QString, 0x00495103,
      20, QMetaType::QString, 0x00495103,
      22, QMetaType::QString, 0x00495103,
      24, QMetaType::QString, 0x00495103,
      26, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       5,
       6,
       7,
       8,
       9,
      10,
      11,
      12,
      13,

       0        // eod
};

void ApproveSignDataCtx687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<ApproveSignDataCtx687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->clicked((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 1: _t->onBack(); break;
        case 2: _t->onApprove(); break;
        case 3: _t->onReject(); break;
        case 4: _t->edited((*reinterpret_cast< QString(*)>(_a[1])),(*reinterpret_cast< QString(*)>(_a[2]))); break;
        case 5: _t->remoteChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 6: _t->transportChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 7: _t->endpointChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 8: _t->fromChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 9: _t->messageChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 10: _t->rawDataChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 11: _t->hashChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 12: _t->passwordChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 13: _t->fromSrcChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (ApproveSignDataCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::clicked)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::onBack)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::onApprove)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::onReject)) {
                *result = 3;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString , QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::edited)) {
                *result = 4;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::remoteChanged)) {
                *result = 5;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::transportChanged)) {
                *result = 6;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::endpointChanged)) {
                *result = 7;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::fromChanged)) {
                *result = 8;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::messageChanged)) {
                *result = 9;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::rawDataChanged)) {
                *result = 10;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::hashChanged)) {
                *result = 11;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::passwordChanged)) {
                *result = 12;
                return;
            }
        }
        {
            using _t = void (ApproveSignDataCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveSignDataCtx687eda::fromSrcChanged)) {
                *result = 13;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<ApproveSignDataCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->remote(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->transport(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->endpoint(); break;
        case 3: *reinterpret_cast< QString*>(_v) = _t->from(); break;
        case 4: *reinterpret_cast< QString*>(_v) = _t->message(); break;
        case 5: *reinterpret_cast< QString*>(_v) = _t->rawData(); break;
        case 6: *reinterpret_cast< QString*>(_v) = _t->hash(); break;
        case 7: *reinterpret_cast< QString*>(_v) = _t->password(); break;
        case 8: *reinterpret_cast< QString*>(_v) = _t->fromSrc(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<ApproveSignDataCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setRemote(*reinterpret_cast< QString*>(_v)); break;
        case 1: _t->setTransport(*reinterpret_cast< QString*>(_v)); break;
        case 2: _t->setEndpoint(*reinterpret_cast< QString*>(_v)); break;
        case 3: _t->setFrom(*reinterpret_cast< QString*>(_v)); break;
        case 4: _t->setMessage(*reinterpret_cast< QString*>(_v)); break;
        case 5: _t->setRawData(*reinterpret_cast< QString*>(_v)); break;
        case 6: _t->setHash(*reinterpret_cast< QString*>(_v)); break;
        case 7: _t->setPassword(*reinterpret_cast< QString*>(_v)); break;
        case 8: _t->setFromSrc(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject ApproveSignDataCtx687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_ApproveSignDataCtx687eda.data,
    qt_meta_data_ApproveSignDataCtx687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *ApproveSignDataCtx687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *ApproveSignDataCtx687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_ApproveSignDataCtx687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int ApproveSignDataCtx687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 14)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 14;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 14)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 14;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 9;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 9;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 9;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 9;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 9;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 9;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void ApproveSignDataCtx687eda::clicked(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}

// SIGNAL 1
void ApproveSignDataCtx687eda::onBack()
{
    QMetaObject::activate(this, &staticMetaObject, 1, nullptr);
}

// SIGNAL 2
void ApproveSignDataCtx687eda::onApprove()
{
    QMetaObject::activate(this, &staticMetaObject, 2, nullptr);
}

// SIGNAL 3
void ApproveSignDataCtx687eda::onReject()
{
    QMetaObject::activate(this, &staticMetaObject, 3, nullptr);
}

// SIGNAL 4
void ApproveSignDataCtx687eda::edited(QString _t1, QString _t2)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))), const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t2))) };
    QMetaObject::activate(this, &staticMetaObject, 4, _a);
}

// SIGNAL 5
void ApproveSignDataCtx687eda::remoteChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 5, _a);
}

// SIGNAL 6
void ApproveSignDataCtx687eda::transportChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 6, _a);
}

// SIGNAL 7
void ApproveSignDataCtx687eda::endpointChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 7, _a);
}

// SIGNAL 8
void ApproveSignDataCtx687eda::fromChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 8, _a);
}

// SIGNAL 9
void ApproveSignDataCtx687eda::messageChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 9, _a);
}

// SIGNAL 10
void ApproveSignDataCtx687eda::rawDataChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 10, _a);
}

// SIGNAL 11
void ApproveSignDataCtx687eda::hashChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 11, _a);
}

// SIGNAL 12
void ApproveSignDataCtx687eda::passwordChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 12, _a);
}

// SIGNAL 13
void ApproveSignDataCtx687eda::fromSrcChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 13, _a);
}
struct qt_meta_stringdata_ApproveTxCtx687eda_t {
    QByteArrayData data[51];
    char stringdata0[556];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_ApproveTxCtx687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_ApproveTxCtx687eda_t qt_meta_stringdata_ApproveTxCtx687eda = {
    {
QT_MOC_LITERAL(0, 0, 18), // "ApproveTxCtx687eda"
QT_MOC_LITERAL(1, 19, 7), // "approve"
QT_MOC_LITERAL(2, 27, 0), // ""
QT_MOC_LITERAL(3, 28, 6), // "reject"
QT_MOC_LITERAL(4, 35, 11), // "checkTxDiff"
QT_MOC_LITERAL(5, 47, 4), // "back"
QT_MOC_LITERAL(6, 52, 6), // "edited"
QT_MOC_LITERAL(7, 59, 1), // "s"
QT_MOC_LITERAL(8, 61, 1), // "v"
QT_MOC_LITERAL(9, 63, 15), // "changeValueUnit"
QT_MOC_LITERAL(10, 79, 18), // "changeGasPriceUnit"
QT_MOC_LITERAL(11, 98, 16), // "valueUnitChanged"
QT_MOC_LITERAL(12, 115, 9), // "valueUnit"
QT_MOC_LITERAL(13, 125, 13), // "remoteChanged"
QT_MOC_LITERAL(14, 139, 6), // "remote"
QT_MOC_LITERAL(15, 146, 16), // "transportChanged"
QT_MOC_LITERAL(16, 163, 9), // "transport"
QT_MOC_LITERAL(17, 173, 15), // "endpointChanged"
QT_MOC_LITERAL(18, 189, 8), // "endpoint"
QT_MOC_LITERAL(19, 198, 11), // "dataChanged"
QT_MOC_LITERAL(20, 210, 4), // "data"
QT_MOC_LITERAL(21, 215, 11), // "fromChanged"
QT_MOC_LITERAL(22, 227, 4), // "from"
QT_MOC_LITERAL(23, 232, 18), // "fromWarningChanged"
QT_MOC_LITERAL(24, 251, 11), // "fromWarning"
QT_MOC_LITERAL(25, 263, 18), // "fromVisibleChanged"
QT_MOC_LITERAL(26, 282, 11), // "fromVisible"
QT_MOC_LITERAL(27, 294, 9), // "toChanged"
QT_MOC_LITERAL(28, 304, 2), // "to"
QT_MOC_LITERAL(29, 307, 16), // "toWarningChanged"
QT_MOC_LITERAL(30, 324, 9), // "toWarning"
QT_MOC_LITERAL(31, 334, 16), // "toVisibleChanged"
QT_MOC_LITERAL(32, 351, 9), // "toVisible"
QT_MOC_LITERAL(33, 361, 10), // "gasChanged"
QT_MOC_LITERAL(34, 372, 3), // "gas"
QT_MOC_LITERAL(35, 376, 15), // "gasPriceChanged"
QT_MOC_LITERAL(36, 392, 8), // "gasPrice"
QT_MOC_LITERAL(37, 401, 19), // "gasPriceUnitChanged"
QT_MOC_LITERAL(38, 421, 12), // "gasPriceUnit"
QT_MOC_LITERAL(39, 434, 12), // "nonceChanged"
QT_MOC_LITERAL(40, 447, 5), // "nonce"
QT_MOC_LITERAL(41, 453, 12), // "valueChanged"
QT_MOC_LITERAL(42, 466, 5), // "value"
QT_MOC_LITERAL(43, 472, 15), // "passwordChanged"
QT_MOC_LITERAL(44, 488, 8), // "password"
QT_MOC_LITERAL(45, 497, 14), // "fromSrcChanged"
QT_MOC_LITERAL(46, 512, 7), // "fromSrc"
QT_MOC_LITERAL(47, 520, 12), // "toSrcChanged"
QT_MOC_LITERAL(48, 533, 5), // "toSrc"
QT_MOC_LITERAL(49, 539, 11), // "diffChanged"
QT_MOC_LITERAL(50, 551, 4) // "diff"

    },
    "ApproveTxCtx687eda\0approve\0\0reject\0"
    "checkTxDiff\0back\0edited\0s\0v\0changeValueUnit\0"
    "changeGasPriceUnit\0valueUnitChanged\0"
    "valueUnit\0remoteChanged\0remote\0"
    "transportChanged\0transport\0endpointChanged\0"
    "endpoint\0dataChanged\0data\0fromChanged\0"
    "from\0fromWarningChanged\0fromWarning\0"
    "fromVisibleChanged\0fromVisible\0toChanged\0"
    "to\0toWarningChanged\0toWarning\0"
    "toVisibleChanged\0toVisible\0gasChanged\0"
    "gas\0gasPriceChanged\0gasPrice\0"
    "gasPriceUnitChanged\0gasPriceUnit\0"
    "nonceChanged\0nonce\0valueChanged\0value\0"
    "passwordChanged\0password\0fromSrcChanged\0"
    "fromSrc\0toSrcChanged\0toSrc\0diffChanged\0"
    "diff"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_ApproveTxCtx687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
      27,   14, // methods
      20,  224, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
      27,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,  149,    2, 0x06 /* Public */,
       3,    0,  150,    2, 0x06 /* Public */,
       4,    0,  151,    2, 0x06 /* Public */,
       5,    0,  152,    2, 0x06 /* Public */,
       6,    2,  153,    2, 0x06 /* Public */,
       9,    1,  158,    2, 0x06 /* Public */,
      10,    1,  161,    2, 0x06 /* Public */,
      11,    1,  164,    2, 0x06 /* Public */,
      13,    1,  167,    2, 0x06 /* Public */,
      15,    1,  170,    2, 0x06 /* Public */,
      17,    1,  173,    2, 0x06 /* Public */,
      19,    1,  176,    2, 0x06 /* Public */,
      21,    1,  179,    2, 0x06 /* Public */,
      23,    1,  182,    2, 0x06 /* Public */,
      25,    1,  185,    2, 0x06 /* Public */,
      27,    1,  188,    2, 0x06 /* Public */,
      29,    1,  191,    2, 0x06 /* Public */,
      31,    1,  194,    2, 0x06 /* Public */,
      33,    1,  197,    2, 0x06 /* Public */,
      35,    1,  200,    2, 0x06 /* Public */,
      37,    1,  203,    2, 0x06 /* Public */,
      39,    1,  206,    2, 0x06 /* Public */,
      41,    1,  209,    2, 0x06 /* Public */,
      43,    1,  212,    2, 0x06 /* Public */,
      45,    1,  215,    2, 0x06 /* Public */,
      47,    1,  218,    2, 0x06 /* Public */,
      49,    1,  221,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void, QMetaType::QString, QMetaType::QString,    7,    8,
    QMetaType::Void, QMetaType::Int,    8,
    QMetaType::Void, QMetaType::Int,    8,
    QMetaType::Void, QMetaType::Int,   12,
    QMetaType::Void, QMetaType::QString,   14,
    QMetaType::Void, QMetaType::QString,   16,
    QMetaType::Void, QMetaType::QString,   18,
    QMetaType::Void, QMetaType::QString,   20,
    QMetaType::Void, QMetaType::QString,   22,
    QMetaType::Void, QMetaType::QString,   24,
    QMetaType::Void, QMetaType::Bool,   26,
    QMetaType::Void, QMetaType::QString,   28,
    QMetaType::Void, QMetaType::QString,   30,
    QMetaType::Void, QMetaType::Bool,   32,
    QMetaType::Void, QMetaType::QString,   34,
    QMetaType::Void, QMetaType::QString,   36,
    QMetaType::Void, QMetaType::Int,   38,
    QMetaType::Void, QMetaType::QString,   40,
    QMetaType::Void, QMetaType::QString,   42,
    QMetaType::Void, QMetaType::QString,   44,
    QMetaType::Void, QMetaType::QString,   46,
    QMetaType::Void, QMetaType::QString,   48,
    QMetaType::Void, QMetaType::QString,   50,

 // properties: name, type, flags
      12, QMetaType::Int, 0x00495103,
      14, QMetaType::QString, 0x00495103,
      16, QMetaType::QString, 0x00495103,
      18, QMetaType::QString, 0x00495103,
      20, QMetaType::QString, 0x00495103,
      22, QMetaType::QString, 0x00495103,
      24, QMetaType::QString, 0x00495103,
      26, QMetaType::Bool, 0x00495103,
      28, QMetaType::QString, 0x00495103,
      30, QMetaType::QString, 0x00495103,
      32, QMetaType::Bool, 0x00495103,
      34, QMetaType::QString, 0x00495103,
      36, QMetaType::QString, 0x00495103,
      38, QMetaType::Int, 0x00495103,
      40, QMetaType::QString, 0x00495103,
      42, QMetaType::QString, 0x00495103,
      44, QMetaType::QString, 0x00495103,
      46, QMetaType::QString, 0x00495103,
      48, QMetaType::QString, 0x00495103,
      50, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       7,
       8,
       9,
      10,
      11,
      12,
      13,
      14,
      15,
      16,
      17,
      18,
      19,
      20,
      21,
      22,
      23,
      24,
      25,
      26,

       0        // eod
};

void ApproveTxCtx687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<ApproveTxCtx687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->approve(); break;
        case 1: _t->reject(); break;
        case 2: _t->checkTxDiff(); break;
        case 3: _t->back(); break;
        case 4: _t->edited((*reinterpret_cast< QString(*)>(_a[1])),(*reinterpret_cast< QString(*)>(_a[2]))); break;
        case 5: _t->changeValueUnit((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 6: _t->changeGasPriceUnit((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 7: _t->valueUnitChanged((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 8: _t->remoteChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 9: _t->transportChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 10: _t->endpointChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 11: _t->dataChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 12: _t->fromChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 13: _t->fromWarningChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 14: _t->fromVisibleChanged((*reinterpret_cast< bool(*)>(_a[1]))); break;
        case 15: _t->toChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 16: _t->toWarningChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 17: _t->toVisibleChanged((*reinterpret_cast< bool(*)>(_a[1]))); break;
        case 18: _t->gasChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 19: _t->gasPriceChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 20: _t->gasPriceUnitChanged((*reinterpret_cast< qint32(*)>(_a[1]))); break;
        case 21: _t->nonceChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 22: _t->valueChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 23: _t->passwordChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 24: _t->fromSrcChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 25: _t->toSrcChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 26: _t->diffChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (ApproveTxCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::approve)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::reject)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::checkTxDiff)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::back)) {
                *result = 3;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString , QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::edited)) {
                *result = 4;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::changeValueUnit)) {
                *result = 5;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::changeGasPriceUnit)) {
                *result = 6;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::valueUnitChanged)) {
                *result = 7;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::remoteChanged)) {
                *result = 8;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::transportChanged)) {
                *result = 9;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::endpointChanged)) {
                *result = 10;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::dataChanged)) {
                *result = 11;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::fromChanged)) {
                *result = 12;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::fromWarningChanged)) {
                *result = 13;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(bool );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::fromVisibleChanged)) {
                *result = 14;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::toChanged)) {
                *result = 15;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::toWarningChanged)) {
                *result = 16;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(bool );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::toVisibleChanged)) {
                *result = 17;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::gasChanged)) {
                *result = 18;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::gasPriceChanged)) {
                *result = 19;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(qint32 );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::gasPriceUnitChanged)) {
                *result = 20;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::nonceChanged)) {
                *result = 21;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::valueChanged)) {
                *result = 22;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::passwordChanged)) {
                *result = 23;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::fromSrcChanged)) {
                *result = 24;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::toSrcChanged)) {
                *result = 25;
                return;
            }
        }
        {
            using _t = void (ApproveTxCtx687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&ApproveTxCtx687eda::diffChanged)) {
                *result = 26;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<ApproveTxCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< qint32*>(_v) = _t->valueUnit(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->remote(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->transport(); break;
        case 3: *reinterpret_cast< QString*>(_v) = _t->endpoint(); break;
        case 4: *reinterpret_cast< QString*>(_v) = _t->data(); break;
        case 5: *reinterpret_cast< QString*>(_v) = _t->from(); break;
        case 6: *reinterpret_cast< QString*>(_v) = _t->fromWarning(); break;
        case 7: *reinterpret_cast< bool*>(_v) = _t->isFromVisible(); break;
        case 8: *reinterpret_cast< QString*>(_v) = _t->to(); break;
        case 9: *reinterpret_cast< QString*>(_v) = _t->toWarning(); break;
        case 10: *reinterpret_cast< bool*>(_v) = _t->isToVisible(); break;
        case 11: *reinterpret_cast< QString*>(_v) = _t->gas(); break;
        case 12: *reinterpret_cast< QString*>(_v) = _t->gasPrice(); break;
        case 13: *reinterpret_cast< qint32*>(_v) = _t->gasPriceUnit(); break;
        case 14: *reinterpret_cast< QString*>(_v) = _t->nonce(); break;
        case 15: *reinterpret_cast< QString*>(_v) = _t->value(); break;
        case 16: *reinterpret_cast< QString*>(_v) = _t->password(); break;
        case 17: *reinterpret_cast< QString*>(_v) = _t->fromSrc(); break;
        case 18: *reinterpret_cast< QString*>(_v) = _t->toSrc(); break;
        case 19: *reinterpret_cast< QString*>(_v) = _t->diff(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<ApproveTxCtx687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setValueUnit(*reinterpret_cast< qint32*>(_v)); break;
        case 1: _t->setRemote(*reinterpret_cast< QString*>(_v)); break;
        case 2: _t->setTransport(*reinterpret_cast< QString*>(_v)); break;
        case 3: _t->setEndpoint(*reinterpret_cast< QString*>(_v)); break;
        case 4: _t->setData(*reinterpret_cast< QString*>(_v)); break;
        case 5: _t->setFrom(*reinterpret_cast< QString*>(_v)); break;
        case 6: _t->setFromWarning(*reinterpret_cast< QString*>(_v)); break;
        case 7: _t->setFromVisible(*reinterpret_cast< bool*>(_v)); break;
        case 8: _t->setTo(*reinterpret_cast< QString*>(_v)); break;
        case 9: _t->setToWarning(*reinterpret_cast< QString*>(_v)); break;
        case 10: _t->setToVisible(*reinterpret_cast< bool*>(_v)); break;
        case 11: _t->setGas(*reinterpret_cast< QString*>(_v)); break;
        case 12: _t->setGasPrice(*reinterpret_cast< QString*>(_v)); break;
        case 13: _t->setGasPriceUnit(*reinterpret_cast< qint32*>(_v)); break;
        case 14: _t->setNonce(*reinterpret_cast< QString*>(_v)); break;
        case 15: _t->setValue(*reinterpret_cast< QString*>(_v)); break;
        case 16: _t->setPassword(*reinterpret_cast< QString*>(_v)); break;
        case 17: _t->setFromSrc(*reinterpret_cast< QString*>(_v)); break;
        case 18: _t->setToSrc(*reinterpret_cast< QString*>(_v)); break;
        case 19: _t->setDiff(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject ApproveTxCtx687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_ApproveTxCtx687eda.data,
    qt_meta_data_ApproveTxCtx687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *ApproveTxCtx687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *ApproveTxCtx687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_ApproveTxCtx687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int ApproveTxCtx687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 27)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 27;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 27)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 27;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 20;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 20;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 20;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 20;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 20;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 20;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void ApproveTxCtx687eda::approve()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void ApproveTxCtx687eda::reject()
{
    QMetaObject::activate(this, &staticMetaObject, 1, nullptr);
}

// SIGNAL 2
void ApproveTxCtx687eda::checkTxDiff()
{
    QMetaObject::activate(this, &staticMetaObject, 2, nullptr);
}

// SIGNAL 3
void ApproveTxCtx687eda::back()
{
    QMetaObject::activate(this, &staticMetaObject, 3, nullptr);
}

// SIGNAL 4
void ApproveTxCtx687eda::edited(QString _t1, QString _t2)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))), const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t2))) };
    QMetaObject::activate(this, &staticMetaObject, 4, _a);
}

// SIGNAL 5
void ApproveTxCtx687eda::changeValueUnit(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 5, _a);
}

// SIGNAL 6
void ApproveTxCtx687eda::changeGasPriceUnit(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 6, _a);
}

// SIGNAL 7
void ApproveTxCtx687eda::valueUnitChanged(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 7, _a);
}

// SIGNAL 8
void ApproveTxCtx687eda::remoteChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 8, _a);
}

// SIGNAL 9
void ApproveTxCtx687eda::transportChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 9, _a);
}

// SIGNAL 10
void ApproveTxCtx687eda::endpointChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 10, _a);
}

// SIGNAL 11
void ApproveTxCtx687eda::dataChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 11, _a);
}

// SIGNAL 12
void ApproveTxCtx687eda::fromChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 12, _a);
}

// SIGNAL 13
void ApproveTxCtx687eda::fromWarningChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 13, _a);
}

// SIGNAL 14
void ApproveTxCtx687eda::fromVisibleChanged(bool _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 14, _a);
}

// SIGNAL 15
void ApproveTxCtx687eda::toChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 15, _a);
}

// SIGNAL 16
void ApproveTxCtx687eda::toWarningChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 16, _a);
}

// SIGNAL 17
void ApproveTxCtx687eda::toVisibleChanged(bool _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 17, _a);
}

// SIGNAL 18
void ApproveTxCtx687eda::gasChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 18, _a);
}

// SIGNAL 19
void ApproveTxCtx687eda::gasPriceChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 19, _a);
}

// SIGNAL 20
void ApproveTxCtx687eda::gasPriceUnitChanged(qint32 _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 20, _a);
}

// SIGNAL 21
void ApproveTxCtx687eda::nonceChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 21, _a);
}

// SIGNAL 22
void ApproveTxCtx687eda::valueChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 22, _a);
}

// SIGNAL 23
void ApproveTxCtx687eda::passwordChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 23, _a);
}

// SIGNAL 24
void ApproveTxCtx687eda::fromSrcChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 24, _a);
}

// SIGNAL 25
void ApproveTxCtx687eda::toSrcChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 25, _a);
}

// SIGNAL 26
void ApproveTxCtx687eda::diffChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 26, _a);
}
struct qt_meta_stringdata_CustomListModel687eda_t {
    QByteArrayData data[9];
    char stringdata0[103];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_CustomListModel687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_CustomListModel687eda_t qt_meta_stringdata_CustomListModel687eda = {
    {
QT_MOC_LITERAL(0, 0, 21), // "CustomListModel687eda"
QT_MOC_LITERAL(1, 22, 5), // "clear"
QT_MOC_LITERAL(2, 28, 0), // ""
QT_MOC_LITERAL(3, 29, 3), // "add"
QT_MOC_LITERAL(4, 33, 8), // "quintptr"
QT_MOC_LITERAL(5, 42, 7), // "account"
QT_MOC_LITERAL(6, 50, 21), // "updateRequiredChanged"
QT_MOC_LITERAL(7, 72, 14), // "updateRequired"
QT_MOC_LITERAL(8, 87, 15) // "callbackFromQml"

    },
    "CustomListModel687eda\0clear\0\0add\0"
    "quintptr\0account\0updateRequiredChanged\0"
    "updateRequired\0callbackFromQml"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_CustomListModel687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       4,   14, // methods
       1,   42, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       3,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   34,    2, 0x06 /* Public */,
       3,    1,   35,    2, 0x06 /* Public */,
       6,    1,   38,    2, 0x06 /* Public */,

 // slots: name, argc, parameters, tag, flags
       8,    0,   41,    2, 0x0a /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void, 0x80000000 | 4,    5,
    QMetaType::Void, QMetaType::QString,    7,

 // slots: parameters
    QMetaType::Void,

 // properties: name, type, flags
       7, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       2,

       0        // eod
};

void CustomListModel687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<CustomListModel687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->clear(); break;
        case 1: _t->add((*reinterpret_cast< quintptr(*)>(_a[1]))); break;
        case 2: _t->updateRequiredChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 3: _t->callbackFromQml(); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (CustomListModel687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&CustomListModel687eda::clear)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (CustomListModel687eda::*)(quintptr );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&CustomListModel687eda::add)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (CustomListModel687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&CustomListModel687eda::updateRequiredChanged)) {
                *result = 2;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<CustomListModel687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->updateRequired(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<CustomListModel687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setUpdateRequired(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject CustomListModel687eda::staticMetaObject = { {
    &QAbstractListModel::staticMetaObject,
    qt_meta_stringdata_CustomListModel687eda.data,
    qt_meta_data_CustomListModel687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *CustomListModel687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *CustomListModel687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_CustomListModel687eda.stringdata0))
        return static_cast<void*>(this);
    return QAbstractListModel::qt_metacast(_clname);
}

int CustomListModel687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QAbstractListModel::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 4)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 4;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 4)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 4;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 1;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void CustomListModel687eda::clear()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void CustomListModel687eda::add(quintptr _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 1, _a);
}

// SIGNAL 2
void CustomListModel687eda::updateRequiredChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}
struct qt_meta_stringdata_LoginContext687eda_t {
    QByteArrayData data[11];
    char stringdata0[111];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_LoginContext687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_LoginContext687eda_t qt_meta_stringdata_LoginContext687eda = {
    {
QT_MOC_LITERAL(0, 0, 18), // "LoginContext687eda"
QT_MOC_LITERAL(1, 19, 5), // "start"
QT_MOC_LITERAL(2, 25, 0), // ""
QT_MOC_LITERAL(3, 26, 9), // "checkPath"
QT_MOC_LITERAL(4, 36, 1), // "b"
QT_MOC_LITERAL(5, 38, 15), // "clefPathChanged"
QT_MOC_LITERAL(6, 54, 8), // "clefPath"
QT_MOC_LITERAL(7, 63, 17), // "binaryHashChanged"
QT_MOC_LITERAL(8, 81, 10), // "binaryHash"
QT_MOC_LITERAL(9, 92, 12), // "errorChanged"
QT_MOC_LITERAL(10, 105, 5) // "error"

    },
    "LoginContext687eda\0start\0\0checkPath\0"
    "b\0clefPathChanged\0clefPath\0binaryHashChanged\0"
    "binaryHash\0errorChanged\0error"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_LoginContext687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       5,   14, // methods
       3,   52, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       5,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   39,    2, 0x06 /* Public */,
       3,    1,   40,    2, 0x06 /* Public */,
       5,    1,   43,    2, 0x06 /* Public */,
       7,    1,   46,    2, 0x06 /* Public */,
       9,    1,   49,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void, QMetaType::QString,    4,
    QMetaType::Void, QMetaType::QString,    6,
    QMetaType::Void, QMetaType::QString,    8,
    QMetaType::Void, QMetaType::QString,   10,

 // properties: name, type, flags
       6, QMetaType::QString, 0x00495103,
       8, QMetaType::QString, 0x00495103,
      10, QMetaType::QString, 0x00495103,

 // properties: notify_signal_id
       2,
       3,
       4,

       0        // eod
};

void LoginContext687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<LoginContext687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->start(); break;
        case 1: _t->checkPath((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 2: _t->clefPathChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 3: _t->binaryHashChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        case 4: _t->errorChanged((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (LoginContext687eda::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&LoginContext687eda::start)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (LoginContext687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&LoginContext687eda::checkPath)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (LoginContext687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&LoginContext687eda::clefPathChanged)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (LoginContext687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&LoginContext687eda::binaryHashChanged)) {
                *result = 3;
                return;
            }
        }
        {
            using _t = void (LoginContext687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&LoginContext687eda::errorChanged)) {
                *result = 4;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<LoginContext687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QString*>(_v) = _t->clefPath(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->binaryHash(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->error(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<LoginContext687eda *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setClefPath(*reinterpret_cast< QString*>(_v)); break;
        case 1: _t->setBinaryHash(*reinterpret_cast< QString*>(_v)); break;
        case 2: _t->setError(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject LoginContext687eda::staticMetaObject = { {
    &QObject::staticMetaObject,
    qt_meta_stringdata_LoginContext687eda.data,
    qt_meta_data_LoginContext687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *LoginContext687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *LoginContext687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_LoginContext687eda.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int LoginContext687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 5)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 5;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 5)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 5;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 3;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void LoginContext687eda::start()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void LoginContext687eda::checkPath(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 1, _a);
}

// SIGNAL 2
void LoginContext687eda::clefPathChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}

// SIGNAL 3
void LoginContext687eda::binaryHashChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 3, _a);
}

// SIGNAL 4
void LoginContext687eda::errorChanged(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 4, _a);
}
struct qt_meta_stringdata_TxListAccountsModel687eda_t {
    QByteArrayData data[4];
    char stringdata0[34];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_TxListAccountsModel687eda_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_TxListAccountsModel687eda_t qt_meta_stringdata_TxListAccountsModel687eda = {
    {
QT_MOC_LITERAL(0, 0, 25), // "TxListAccountsModel687eda"
QT_MOC_LITERAL(1, 26, 3), // "add"
QT_MOC_LITERAL(2, 30, 0), // ""
QT_MOC_LITERAL(3, 31, 2) // "tx"

    },
    "TxListAccountsModel687eda\0add\0\0tx"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_TxListAccountsModel687eda[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
       1,   14, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       1,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   19,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, QMetaType::QString,    3,

       0        // eod
};

void TxListAccountsModel687eda::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<TxListAccountsModel687eda *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->add((*reinterpret_cast< QString(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (TxListAccountsModel687eda::*)(QString );
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&TxListAccountsModel687eda::add)) {
                *result = 0;
                return;
            }
        }
    }
}

QT_INIT_METAOBJECT const QMetaObject TxListAccountsModel687eda::staticMetaObject = { {
    &QAbstractListModel::staticMetaObject,
    qt_meta_stringdata_TxListAccountsModel687eda.data,
    qt_meta_data_TxListAccountsModel687eda,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *TxListAccountsModel687eda::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *TxListAccountsModel687eda::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_TxListAccountsModel687eda.stringdata0))
        return static_cast<void*>(this);
    return QAbstractListModel::qt_metacast(_clname);
}

int TxListAccountsModel687eda::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QAbstractListModel::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 1)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 1;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 1)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 1;
    }
    return _id;
}

// SIGNAL 0
void TxListAccountsModel687eda::add(QString _t1)
{
    void *_a[] = { nullptr, const_cast<void*>(reinterpret_cast<const void*>(std::addressof(_t1))) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}
QT_WARNING_POP
QT_END_MOC_NAMESPACE
