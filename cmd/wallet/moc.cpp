

#define protected public
#define private public

#include "moc.h"
#include "_cgo_export.h"

#include <QAbstractItemModel>
#include <QAbstractListModel>
#include <QByteArray>
#include <QChildEvent>
#include <QEvent>
#include <QGraphicsObject>
#include <QGraphicsWidget>
#include <QHash>
#include <QLayout>
#include <QMap>
#include <QMetaMethod>
#include <QMetaObject>
#include <QMimeData>
#include <QModelIndex>
#include <QObject>
#include <QOffscreenSurface>
#include <QPaintDeviceWindow>
#include <QPdfWriter>
#include <QPersistentModelIndex>
#include <QQuickItem>
#include <QSize>
#include <QString>
#include <QTimerEvent>
#include <QVariant>
#include <QVector>
#include <QWidget>
#include <QWindow>

#ifdef QT_QML_LIB
	#include <QQmlEngine>
#endif


class ApproveListingCtx687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(QString remote READ remote WRITE setRemote NOTIFY remoteChanged)
Q_PROPERTY(QString transport READ transport WRITE setTransport NOTIFY transportChanged)
Q_PROPERTY(QString endpoint READ endpoint WRITE setEndpoint NOTIFY endpointChanged)
Q_PROPERTY(QString from READ from WRITE setFrom NOTIFY fromChanged)
Q_PROPERTY(QString message READ message WRITE setMessage NOTIFY messageChanged)
Q_PROPERTY(QString rawData READ rawData WRITE setRawData NOTIFY rawDataChanged)
Q_PROPERTY(QString hash READ hash WRITE setHash NOTIFY hashChanged)
public:
	ApproveListingCtx687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType();ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaTypes();callbackApproveListingCtx687eda_Constructor(this);};
	void Signal_Back() { callbackApproveListingCtx687eda_Back(this); };
	void Signal_Approve() { callbackApproveListingCtx687eda_Approve(this); };
	void Signal_Reject() { callbackApproveListingCtx687eda_Reject(this); };
	void Signal_OnCheckStateChanged(qint32 i, bool checked) { callbackApproveListingCtx687eda_OnCheckStateChanged(this, i, checked); };
	void Signal_TriggerUpdate() { callbackApproveListingCtx687eda_TriggerUpdate(this); };
	QString remote() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_Remote(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRemote(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveListingCtx687eda_SetRemote(this, remotePacked); };
	void Signal_RemoteChanged(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveListingCtx687eda_RemoteChanged(this, remotePacked); };
	QString transport() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_Transport(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setTransport(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveListingCtx687eda_SetTransport(this, transportPacked); };
	void Signal_TransportChanged(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveListingCtx687eda_TransportChanged(this, transportPacked); };
	QString endpoint() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_Endpoint(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setEndpoint(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveListingCtx687eda_SetEndpoint(this, endpointPacked); };
	void Signal_EndpointChanged(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveListingCtx687eda_EndpointChanged(this, endpointPacked); };
	QString from() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_From(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFrom(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveListingCtx687eda_SetFrom(this, fromPacked); };
	void Signal_FromChanged(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveListingCtx687eda_FromChanged(this, fromPacked); };
	QString message() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_Message(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setMessage(QString message) { QByteArray* t6f9b9a = new QByteArray(message.toUtf8()); Moc_PackedString messagePacked = { const_cast<char*>(t6f9b9a->prepend("WHITESPACE").constData()+10), t6f9b9a->size()-10, t6f9b9a };callbackApproveListingCtx687eda_SetMessage(this, messagePacked); };
	void Signal_MessageChanged(QString message) { QByteArray* t6f9b9a = new QByteArray(message.toUtf8()); Moc_PackedString messagePacked = { const_cast<char*>(t6f9b9a->prepend("WHITESPACE").constData()+10), t6f9b9a->size()-10, t6f9b9a };callbackApproveListingCtx687eda_MessageChanged(this, messagePacked); };
	QString rawData() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_RawData(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRawData(QString rawData) { QByteArray* tcacc10 = new QByteArray(rawData.toUtf8()); Moc_PackedString rawDataPacked = { const_cast<char*>(tcacc10->prepend("WHITESPACE").constData()+10), tcacc10->size()-10, tcacc10 };callbackApproveListingCtx687eda_SetRawData(this, rawDataPacked); };
	void Signal_RawDataChanged(QString rawData) { QByteArray* tcacc10 = new QByteArray(rawData.toUtf8()); Moc_PackedString rawDataPacked = { const_cast<char*>(tcacc10->prepend("WHITESPACE").constData()+10), tcacc10->size()-10, tcacc10 };callbackApproveListingCtx687eda_RawDataChanged(this, rawDataPacked); };
	QString hash() { return ({ Moc_PackedString tempVal = callbackApproveListingCtx687eda_Hash(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setHash(QString hash) { QByteArray* t2346ad = new QByteArray(hash.toUtf8()); Moc_PackedString hashPacked = { const_cast<char*>(t2346ad->prepend("WHITESPACE").constData()+10), t2346ad->size()-10, t2346ad };callbackApproveListingCtx687eda_SetHash(this, hashPacked); };
	void Signal_HashChanged(QString hash) { QByteArray* t2346ad = new QByteArray(hash.toUtf8()); Moc_PackedString hashPacked = { const_cast<char*>(t2346ad->prepend("WHITESPACE").constData()+10), t2346ad->size()-10, t2346ad };callbackApproveListingCtx687eda_HashChanged(this, hashPacked); };
	 ~ApproveListingCtx687eda() { callbackApproveListingCtx687eda_DestroyApproveListingCtx(this); };
	void childEvent(QChildEvent * event) { callbackApproveListingCtx687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackApproveListingCtx687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackApproveListingCtx687eda_CustomEvent(this, event); };
	void deleteLater() { callbackApproveListingCtx687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackApproveListingCtx687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackApproveListingCtx687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackApproveListingCtx687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackApproveListingCtx687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackApproveListingCtx687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackApproveListingCtx687eda_TimerEvent(this, event); };
	QString remoteDefault() { return _remote; };
	void setRemoteDefault(QString p) { if (p != _remote) { _remote = p; remoteChanged(_remote); } };
	QString transportDefault() { return _transport; };
	void setTransportDefault(QString p) { if (p != _transport) { _transport = p; transportChanged(_transport); } };
	QString endpointDefault() { return _endpoint; };
	void setEndpointDefault(QString p) { if (p != _endpoint) { _endpoint = p; endpointChanged(_endpoint); } };
	QString fromDefault() { return _from; };
	void setFromDefault(QString p) { if (p != _from) { _from = p; fromChanged(_from); } };
	QString messageDefault() { return _message; };
	void setMessageDefault(QString p) { if (p != _message) { _message = p; messageChanged(_message); } };
	QString rawDataDefault() { return _rawData; };
	void setRawDataDefault(QString p) { if (p != _rawData) { _rawData = p; rawDataChanged(_rawData); } };
	QString hashDefault() { return _hash; };
	void setHashDefault(QString p) { if (p != _hash) { _hash = p; hashChanged(_hash); } };
signals:
	void back();
	void approve();
	void reject();
	void onCheckStateChanged(qint32 i, bool checked);
	void triggerUpdate();
	void remoteChanged(QString remote);
	void transportChanged(QString transport);
	void endpointChanged(QString endpoint);
	void fromChanged(QString from);
	void messageChanged(QString message);
	void rawDataChanged(QString rawData);
	void hashChanged(QString hash);
public slots:
private:
	QString _remote;
	QString _transport;
	QString _endpoint;
	QString _from;
	QString _message;
	QString _rawData;
	QString _hash;
};

Q_DECLARE_METATYPE(ApproveListingCtx687eda*)


void ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class ApproveNewAccountCtx687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(QString remote READ remote WRITE setRemote NOTIFY remoteChanged)
Q_PROPERTY(QString transport READ transport WRITE setTransport NOTIFY transportChanged)
Q_PROPERTY(QString endpoint READ endpoint WRITE setEndpoint NOTIFY endpointChanged)
Q_PROPERTY(QString password READ password WRITE setPassword NOTIFY passwordChanged)
Q_PROPERTY(QString confirmPassword READ confirmPassword WRITE setConfirmPassword NOTIFY confirmPasswordChanged)
public:
	ApproveNewAccountCtx687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType();ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaTypes();callbackApproveNewAccountCtx687eda_Constructor(this);};
	void Signal_Clicked(qint32 b) { callbackApproveNewAccountCtx687eda_Clicked(this, b); };
	void Signal_Back() { callbackApproveNewAccountCtx687eda_Back(this); };
	void Signal_PasswordEdited(QString b) { QByteArray* te9d71f = new QByteArray(b.toUtf8()); Moc_PackedString bPacked = { const_cast<char*>(te9d71f->prepend("WHITESPACE").constData()+10), te9d71f->size()-10, te9d71f };callbackApproveNewAccountCtx687eda_PasswordEdited(this, bPacked); };
	void Signal_ConfirmPasswordEdited(QString b) { QByteArray* te9d71f = new QByteArray(b.toUtf8()); Moc_PackedString bPacked = { const_cast<char*>(te9d71f->prepend("WHITESPACE").constData()+10), te9d71f->size()-10, te9d71f };callbackApproveNewAccountCtx687eda_ConfirmPasswordEdited(this, bPacked); };
	QString remote() { return ({ Moc_PackedString tempVal = callbackApproveNewAccountCtx687eda_Remote(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRemote(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveNewAccountCtx687eda_SetRemote(this, remotePacked); };
	void Signal_RemoteChanged(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveNewAccountCtx687eda_RemoteChanged(this, remotePacked); };
	QString transport() { return ({ Moc_PackedString tempVal = callbackApproveNewAccountCtx687eda_Transport(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setTransport(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveNewAccountCtx687eda_SetTransport(this, transportPacked); };
	void Signal_TransportChanged(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveNewAccountCtx687eda_TransportChanged(this, transportPacked); };
	QString endpoint() { return ({ Moc_PackedString tempVal = callbackApproveNewAccountCtx687eda_Endpoint(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setEndpoint(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveNewAccountCtx687eda_SetEndpoint(this, endpointPacked); };
	void Signal_EndpointChanged(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveNewAccountCtx687eda_EndpointChanged(this, endpointPacked); };
	QString password() { return ({ Moc_PackedString tempVal = callbackApproveNewAccountCtx687eda_Password(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setPassword(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveNewAccountCtx687eda_SetPassword(this, passwordPacked); };
	void Signal_PasswordChanged(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveNewAccountCtx687eda_PasswordChanged(this, passwordPacked); };
	QString confirmPassword() { return ({ Moc_PackedString tempVal = callbackApproveNewAccountCtx687eda_ConfirmPassword(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setConfirmPassword(QString confirmPassword) { QByteArray* t802b53 = new QByteArray(confirmPassword.toUtf8()); Moc_PackedString confirmPasswordPacked = { const_cast<char*>(t802b53->prepend("WHITESPACE").constData()+10), t802b53->size()-10, t802b53 };callbackApproveNewAccountCtx687eda_SetConfirmPassword(this, confirmPasswordPacked); };
	void Signal_ConfirmPasswordChanged(QString confirmPassword) { QByteArray* t802b53 = new QByteArray(confirmPassword.toUtf8()); Moc_PackedString confirmPasswordPacked = { const_cast<char*>(t802b53->prepend("WHITESPACE").constData()+10), t802b53->size()-10, t802b53 };callbackApproveNewAccountCtx687eda_ConfirmPasswordChanged(this, confirmPasswordPacked); };
	 ~ApproveNewAccountCtx687eda() { callbackApproveNewAccountCtx687eda_DestroyApproveNewAccountCtx(this); };
	void childEvent(QChildEvent * event) { callbackApproveNewAccountCtx687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackApproveNewAccountCtx687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackApproveNewAccountCtx687eda_CustomEvent(this, event); };
	void deleteLater() { callbackApproveNewAccountCtx687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackApproveNewAccountCtx687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackApproveNewAccountCtx687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackApproveNewAccountCtx687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackApproveNewAccountCtx687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackApproveNewAccountCtx687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackApproveNewAccountCtx687eda_TimerEvent(this, event); };
	QString remoteDefault() { return _remote; };
	void setRemoteDefault(QString p) { if (p != _remote) { _remote = p; remoteChanged(_remote); } };
	QString transportDefault() { return _transport; };
	void setTransportDefault(QString p) { if (p != _transport) { _transport = p; transportChanged(_transport); } };
	QString endpointDefault() { return _endpoint; };
	void setEndpointDefault(QString p) { if (p != _endpoint) { _endpoint = p; endpointChanged(_endpoint); } };
	QString passwordDefault() { return _password; };
	void setPasswordDefault(QString p) { if (p != _password) { _password = p; passwordChanged(_password); } };
	QString confirmPasswordDefault() { return _confirmPassword; };
	void setConfirmPasswordDefault(QString p) { if (p != _confirmPassword) { _confirmPassword = p; confirmPasswordChanged(_confirmPassword); } };
signals:
	void clicked(qint32 b);
	void back();
	void passwordEdited(QString b);
	void confirmPasswordEdited(QString b);
	void remoteChanged(QString remote);
	void transportChanged(QString transport);
	void endpointChanged(QString endpoint);
	void passwordChanged(QString password);
	void confirmPasswordChanged(QString confirmPassword);
public slots:
private:
	QString _remote;
	QString _transport;
	QString _endpoint;
	QString _password;
	QString _confirmPassword;
};

Q_DECLARE_METATYPE(ApproveNewAccountCtx687eda*)


void ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class TxListCtx687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(QString shortenAddress READ shortenAddress WRITE setShortenAddress NOTIFY shortenAddressChanged)
Q_PROPERTY(QString selectedSrc READ selectedSrc WRITE setSelectedSrc NOTIFY selectedSrcChanged)
public:
	TxListCtx687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");TxListCtx687eda_TxListCtx687eda_QRegisterMetaType();TxListCtx687eda_TxListCtx687eda_QRegisterMetaTypes();callbackTxListCtx687eda_Constructor(this);};
	void Signal_Clicked(qint32 b) { callbackTxListCtx687eda_Clicked(this, b); };
	QString shortenAddress() { return ({ Moc_PackedString tempVal = callbackTxListCtx687eda_ShortenAddress(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setShortenAddress(QString shortenAddress) { QByteArray* t3fdf7b = new QByteArray(shortenAddress.toUtf8()); Moc_PackedString shortenAddressPacked = { const_cast<char*>(t3fdf7b->prepend("WHITESPACE").constData()+10), t3fdf7b->size()-10, t3fdf7b };callbackTxListCtx687eda_SetShortenAddress(this, shortenAddressPacked); };
	void Signal_ShortenAddressChanged(QString shortenAddress) { QByteArray* t3fdf7b = new QByteArray(shortenAddress.toUtf8()); Moc_PackedString shortenAddressPacked = { const_cast<char*>(t3fdf7b->prepend("WHITESPACE").constData()+10), t3fdf7b->size()-10, t3fdf7b };callbackTxListCtx687eda_ShortenAddressChanged(this, shortenAddressPacked); };
	QString selectedSrc() { return ({ Moc_PackedString tempVal = callbackTxListCtx687eda_SelectedSrc(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setSelectedSrc(QString selectedSrc) { QByteArray* t5f7742 = new QByteArray(selectedSrc.toUtf8()); Moc_PackedString selectedSrcPacked = { const_cast<char*>(t5f7742->prepend("WHITESPACE").constData()+10), t5f7742->size()-10, t5f7742 };callbackTxListCtx687eda_SetSelectedSrc(this, selectedSrcPacked); };
	void Signal_SelectedSrcChanged(QString selectedSrc) { QByteArray* t5f7742 = new QByteArray(selectedSrc.toUtf8()); Moc_PackedString selectedSrcPacked = { const_cast<char*>(t5f7742->prepend("WHITESPACE").constData()+10), t5f7742->size()-10, t5f7742 };callbackTxListCtx687eda_SelectedSrcChanged(this, selectedSrcPacked); };
	 ~TxListCtx687eda() { callbackTxListCtx687eda_DestroyTxListCtx(this); };
	void childEvent(QChildEvent * event) { callbackTxListCtx687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackTxListCtx687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackTxListCtx687eda_CustomEvent(this, event); };
	void deleteLater() { callbackTxListCtx687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackTxListCtx687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackTxListCtx687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackTxListCtx687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackTxListCtx687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackTxListCtx687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackTxListCtx687eda_TimerEvent(this, event); };
	QString shortenAddressDefault() { return _shortenAddress; };
	void setShortenAddressDefault(QString p) { if (p != _shortenAddress) { _shortenAddress = p; shortenAddressChanged(_shortenAddress); } };
	QString selectedSrcDefault() { return _selectedSrc; };
	void setSelectedSrcDefault(QString p) { if (p != _selectedSrc) { _selectedSrc = p; selectedSrcChanged(_selectedSrc); } };
signals:
	void clicked(qint32 b);
	void shortenAddressChanged(QString shortenAddress);
	void selectedSrcChanged(QString selectedSrc);
public slots:
private:
	QString _shortenAddress;
	QString _selectedSrc;
};

Q_DECLARE_METATYPE(TxListCtx687eda*)


void TxListCtx687eda_TxListCtx687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class TxListModel687eda: public QAbstractListModel
{
Q_OBJECT
Q_PROPERTY(bool isEmpty READ isEmpty WRITE setIsEmpty NOTIFY isEmptyChanged)
public:
	TxListModel687eda(QObject *parent = Q_NULLPTR) : QAbstractListModel(parent) {qRegisterMetaType<quintptr>("quintptr");TxListModel687eda_TxListModel687eda_QRegisterMetaType();TxListModel687eda_TxListModel687eda_QRegisterMetaTypes();callbackTxListModel687eda_Constructor(this);};
	void Signal_Clear() { callbackTxListModel687eda_Clear(this); };
	void Signal_Add(quintptr tx) { callbackTxListModel687eda_Add(this, tx); };
	void Signal_Remove(qint32 i) { callbackTxListModel687eda_Remove(this, i); };
	bool isEmpty() { return callbackTxListModel687eda_IsEmpty(this) != 0; };
	void setIsEmpty(bool isEmpty) { callbackTxListModel687eda_SetIsEmpty(this, isEmpty); };
	void Signal_IsEmptyChanged(bool isEmpty) { callbackTxListModel687eda_IsEmptyChanged(this, isEmpty); };
	 ~TxListModel687eda() { callbackTxListModel687eda_DestroyTxListModel(this); };
	bool dropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) { return callbackTxListModel687eda_DropMimeData(this, const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	Qt::ItemFlags flags(const QModelIndex & index) const { return static_cast<Qt::ItemFlag>(callbackTxListModel687eda_Flags(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	QModelIndex index(int row, int column, const QModelIndex & parent) const { return *static_cast<QModelIndex*>(callbackTxListModel687eda_Index(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&parent))); };
	QModelIndex sibling(int row, int column, const QModelIndex & idx) const { return *static_cast<QModelIndex*>(callbackTxListModel687eda_Sibling(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&idx))); };
	QModelIndex buddy(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackTxListModel687eda_Buddy(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool canDropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) const { return callbackTxListModel687eda_CanDropMimeData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	bool canFetchMore(const QModelIndex & parent) const { return callbackTxListModel687eda_CanFetchMore(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	int columnCount(const QModelIndex & parent) const { return callbackTxListModel687eda_ColumnCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_ColumnsAboutToBeInserted(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_ColumnsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationColumn) { callbackTxListModel687eda_ColumnsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationColumn); };
	void Signal_ColumnsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_ColumnsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsInserted(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_ColumnsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int column) { callbackTxListModel687eda_ColumnsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), column); };
	void Signal_ColumnsRemoved(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_ColumnsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	QVariant data(const QModelIndex & index, int role) const { return *static_cast<QVariant*>(callbackTxListModel687eda_Data(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index), role)); };
	void Signal_DataChanged(const QModelIndex & topLeft, const QModelIndex & bottomRight, const QVector<int> & roles) { callbackTxListModel687eda_DataChanged(this, const_cast<QModelIndex*>(&topLeft), const_cast<QModelIndex*>(&bottomRight), ({ QVector<int>* tmpValue037c88 = new QVector<int>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })); };
	void fetchMore(const QModelIndex & parent) { callbackTxListModel687eda_FetchMore(this, const_cast<QModelIndex*>(&parent)); };
	bool hasChildren(const QModelIndex & parent) const { return callbackTxListModel687eda_HasChildren(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	QVariant headerData(int section, Qt::Orientation orientation, int role) const { return *static_cast<QVariant*>(callbackTxListModel687eda_HeaderData(const_cast<void*>(static_cast<const void*>(this)), section, orientation, role)); };
	void Signal_HeaderDataChanged(Qt::Orientation orientation, int first, int last) { callbackTxListModel687eda_HeaderDataChanged(this, orientation, first, last); };
	bool insertColumns(int column, int count, const QModelIndex & parent) { return callbackTxListModel687eda_InsertColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool insertRows(int row, int count, const QModelIndex & parent) { return callbackTxListModel687eda_InsertRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	QMap<int, QVariant> itemData(const QModelIndex & index) const { return ({ QMap<int, QVariant>* tmpP = static_cast<QMap<int, QVariant>*>(callbackTxListModel687eda_ItemData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); QMap<int, QVariant> tmpV = *tmpP; tmpP->~QMap(); free(tmpP); tmpV; }); };
	void Signal_LayoutAboutToBeChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackTxListModel687eda_LayoutAboutToBeChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	void Signal_LayoutChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackTxListModel687eda_LayoutChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	QList<QModelIndex> match(const QModelIndex & start, int role, const QVariant & value, int hits, Qt::MatchFlags flags) const { return ({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(callbackTxListModel687eda_Match(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&start), role, const_cast<QVariant*>(&value), hits, flags)); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }); };
	QMimeData * mimeData(const QModelIndexList & indexes) const { return static_cast<QMimeData*>(callbackTxListModel687eda_MimeData(const_cast<void*>(static_cast<const void*>(this)), ({ QList<QModelIndex>* tmpValuee0adf2 = new QList<QModelIndex>(indexes); Moc_PackedList { tmpValuee0adf2, tmpValuee0adf2->size() }; }))); };
	QStringList mimeTypes() const { return ({ Moc_PackedString tempVal = callbackTxListModel687eda_MimeTypes(const_cast<void*>(static_cast<const void*>(this))); QStringList ret = QString::fromUtf8(tempVal.data, tempVal.len).split("¡¦!", QString::SkipEmptyParts); free(tempVal.data); ret; }); };
	void Signal_ModelAboutToBeReset() { callbackTxListModel687eda_ModelAboutToBeReset(this); };
	void Signal_ModelReset() { callbackTxListModel687eda_ModelReset(this); };
	bool moveColumns(const QModelIndex & sourceParent, int sourceColumn, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackTxListModel687eda_MoveColumns(this, const_cast<QModelIndex*>(&sourceParent), sourceColumn, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	bool moveRows(const QModelIndex & sourceParent, int sourceRow, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackTxListModel687eda_MoveRows(this, const_cast<QModelIndex*>(&sourceParent), sourceRow, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	QModelIndex parent(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackTxListModel687eda_Parent(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool removeColumns(int column, int count, const QModelIndex & parent) { return callbackTxListModel687eda_RemoveColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool removeRows(int row, int count, const QModelIndex & parent) { return callbackTxListModel687eda_RemoveRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	void resetInternalData() { callbackTxListModel687eda_ResetInternalData(this); };
	void revert() { callbackTxListModel687eda_Revert(this); };
	QHash<int, QByteArray> roleNames() const { return ({ QHash<int, QByteArray>* tmpP = static_cast<QHash<int, QByteArray>*>(callbackTxListModel687eda_RoleNames(const_cast<void*>(static_cast<const void*>(this)))); QHash<int, QByteArray> tmpV = *tmpP; tmpP->~QHash(); free(tmpP); tmpV; }); };
	int rowCount(const QModelIndex & parent) const { return callbackTxListModel687eda_RowCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_RowsAboutToBeInserted(const QModelIndex & parent, int start, int end) { callbackTxListModel687eda_RowsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), start, end); };
	void Signal_RowsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationRow) { callbackTxListModel687eda_RowsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationRow); };
	void Signal_RowsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_RowsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsInserted(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_RowsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int row) { callbackTxListModel687eda_RowsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), row); };
	void Signal_RowsRemoved(const QModelIndex & parent, int first, int last) { callbackTxListModel687eda_RowsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	bool setData(const QModelIndex & index, const QVariant & value, int role) { return callbackTxListModel687eda_SetData(this, const_cast<QModelIndex*>(&index), const_cast<QVariant*>(&value), role) != 0; };
	bool setHeaderData(int section, Qt::Orientation orientation, const QVariant & value, int role) { return callbackTxListModel687eda_SetHeaderData(this, section, orientation, const_cast<QVariant*>(&value), role) != 0; };
	bool setItemData(const QModelIndex & index, const QMap<int, QVariant> & roles) { return callbackTxListModel687eda_SetItemData(this, const_cast<QModelIndex*>(&index), ({ QMap<int, QVariant>* tmpValue037c88 = new QMap<int, QVariant>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })) != 0; };
	void sort(int column, Qt::SortOrder order) { callbackTxListModel687eda_Sort(this, column, order); };
	QSize span(const QModelIndex & index) const { return *static_cast<QSize*>(callbackTxListModel687eda_Span(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool submit() { return callbackTxListModel687eda_Submit(this) != 0; };
	Qt::DropActions supportedDragActions() const { return static_cast<Qt::DropAction>(callbackTxListModel687eda_SupportedDragActions(const_cast<void*>(static_cast<const void*>(this)))); };
	Qt::DropActions supportedDropActions() const { return static_cast<Qt::DropAction>(callbackTxListModel687eda_SupportedDropActions(const_cast<void*>(static_cast<const void*>(this)))); };
	void childEvent(QChildEvent * event) { callbackTxListModel687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackTxListModel687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackTxListModel687eda_CustomEvent(this, event); };
	void deleteLater() { callbackTxListModel687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackTxListModel687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackTxListModel687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackTxListModel687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackTxListModel687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackTxListModel687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackTxListModel687eda_TimerEvent(this, event); };
	bool isEmptyDefault() { return _isEmpty; };
	void setIsEmptyDefault(bool p) { if (p != _isEmpty) { _isEmpty = p; isEmptyChanged(_isEmpty); } };
signals:
	void clear();
	void add(quintptr tx);
	void remove(qint32 i);
	void isEmptyChanged(bool isEmpty);
public slots:
private:
	bool _isEmpty;
};

Q_DECLARE_METATYPE(TxListModel687eda*)


void TxListModel687eda_TxListModel687eda_QRegisterMetaTypes() {
}

class ApproveSignDataCtx687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(QString remote READ remote WRITE setRemote NOTIFY remoteChanged)
Q_PROPERTY(QString transport READ transport WRITE setTransport NOTIFY transportChanged)
Q_PROPERTY(QString endpoint READ endpoint WRITE setEndpoint NOTIFY endpointChanged)
Q_PROPERTY(QString from READ from WRITE setFrom NOTIFY fromChanged)
Q_PROPERTY(QString message READ message WRITE setMessage NOTIFY messageChanged)
Q_PROPERTY(QString rawData READ rawData WRITE setRawData NOTIFY rawDataChanged)
Q_PROPERTY(QString hash READ hash WRITE setHash NOTIFY hashChanged)
Q_PROPERTY(QString password READ password WRITE setPassword NOTIFY passwordChanged)
Q_PROPERTY(QString fromSrc READ fromSrc WRITE setFromSrc NOTIFY fromSrcChanged)
public:
	ApproveSignDataCtx687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType();ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaTypes();callbackApproveSignDataCtx687eda_Constructor(this);};
	void Signal_Clicked(qint32 b) { callbackApproveSignDataCtx687eda_Clicked(this, b); };
	void Signal_OnBack() { callbackApproveSignDataCtx687eda_OnBack(this); };
	void Signal_OnApprove() { callbackApproveSignDataCtx687eda_OnApprove(this); };
	void Signal_OnReject() { callbackApproveSignDataCtx687eda_OnReject(this); };
	void Signal_Edited(QString b, QString value) { QByteArray* te9d71f = new QByteArray(b.toUtf8()); Moc_PackedString bPacked = { const_cast<char*>(te9d71f->prepend("WHITESPACE").constData()+10), te9d71f->size()-10, te9d71f };QByteArray* tf32b67 = new QByteArray(value.toUtf8()); Moc_PackedString valuePacked = { const_cast<char*>(tf32b67->prepend("WHITESPACE").constData()+10), tf32b67->size()-10, tf32b67 };callbackApproveSignDataCtx687eda_Edited(this, bPacked, valuePacked); };
	QString remote() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Remote(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRemote(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveSignDataCtx687eda_SetRemote(this, remotePacked); };
	void Signal_RemoteChanged(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveSignDataCtx687eda_RemoteChanged(this, remotePacked); };
	QString transport() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Transport(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setTransport(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveSignDataCtx687eda_SetTransport(this, transportPacked); };
	void Signal_TransportChanged(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveSignDataCtx687eda_TransportChanged(this, transportPacked); };
	QString endpoint() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Endpoint(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setEndpoint(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveSignDataCtx687eda_SetEndpoint(this, endpointPacked); };
	void Signal_EndpointChanged(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveSignDataCtx687eda_EndpointChanged(this, endpointPacked); };
	QString from() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_From(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFrom(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveSignDataCtx687eda_SetFrom(this, fromPacked); };
	void Signal_FromChanged(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveSignDataCtx687eda_FromChanged(this, fromPacked); };
	QString message() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Message(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setMessage(QString message) { QByteArray* t6f9b9a = new QByteArray(message.toUtf8()); Moc_PackedString messagePacked = { const_cast<char*>(t6f9b9a->prepend("WHITESPACE").constData()+10), t6f9b9a->size()-10, t6f9b9a };callbackApproveSignDataCtx687eda_SetMessage(this, messagePacked); };
	void Signal_MessageChanged(QString message) { QByteArray* t6f9b9a = new QByteArray(message.toUtf8()); Moc_PackedString messagePacked = { const_cast<char*>(t6f9b9a->prepend("WHITESPACE").constData()+10), t6f9b9a->size()-10, t6f9b9a };callbackApproveSignDataCtx687eda_MessageChanged(this, messagePacked); };
	QString rawData() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_RawData(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRawData(QString rawData) { QByteArray* tcacc10 = new QByteArray(rawData.toUtf8()); Moc_PackedString rawDataPacked = { const_cast<char*>(tcacc10->prepend("WHITESPACE").constData()+10), tcacc10->size()-10, tcacc10 };callbackApproveSignDataCtx687eda_SetRawData(this, rawDataPacked); };
	void Signal_RawDataChanged(QString rawData) { QByteArray* tcacc10 = new QByteArray(rawData.toUtf8()); Moc_PackedString rawDataPacked = { const_cast<char*>(tcacc10->prepend("WHITESPACE").constData()+10), tcacc10->size()-10, tcacc10 };callbackApproveSignDataCtx687eda_RawDataChanged(this, rawDataPacked); };
	QString hash() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Hash(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setHash(QString hash) { QByteArray* t2346ad = new QByteArray(hash.toUtf8()); Moc_PackedString hashPacked = { const_cast<char*>(t2346ad->prepend("WHITESPACE").constData()+10), t2346ad->size()-10, t2346ad };callbackApproveSignDataCtx687eda_SetHash(this, hashPacked); };
	void Signal_HashChanged(QString hash) { QByteArray* t2346ad = new QByteArray(hash.toUtf8()); Moc_PackedString hashPacked = { const_cast<char*>(t2346ad->prepend("WHITESPACE").constData()+10), t2346ad->size()-10, t2346ad };callbackApproveSignDataCtx687eda_HashChanged(this, hashPacked); };
	QString password() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_Password(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setPassword(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveSignDataCtx687eda_SetPassword(this, passwordPacked); };
	void Signal_PasswordChanged(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveSignDataCtx687eda_PasswordChanged(this, passwordPacked); };
	QString fromSrc() { return ({ Moc_PackedString tempVal = callbackApproveSignDataCtx687eda_FromSrc(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFromSrc(QString fromSrc) { QByteArray* ta8ded4 = new QByteArray(fromSrc.toUtf8()); Moc_PackedString fromSrcPacked = { const_cast<char*>(ta8ded4->prepend("WHITESPACE").constData()+10), ta8ded4->size()-10, ta8ded4 };callbackApproveSignDataCtx687eda_SetFromSrc(this, fromSrcPacked); };
	void Signal_FromSrcChanged(QString fromSrc) { QByteArray* ta8ded4 = new QByteArray(fromSrc.toUtf8()); Moc_PackedString fromSrcPacked = { const_cast<char*>(ta8ded4->prepend("WHITESPACE").constData()+10), ta8ded4->size()-10, ta8ded4 };callbackApproveSignDataCtx687eda_FromSrcChanged(this, fromSrcPacked); };
	 ~ApproveSignDataCtx687eda() { callbackApproveSignDataCtx687eda_DestroyApproveSignDataCtx(this); };
	void childEvent(QChildEvent * event) { callbackApproveSignDataCtx687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackApproveSignDataCtx687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackApproveSignDataCtx687eda_CustomEvent(this, event); };
	void deleteLater() { callbackApproveSignDataCtx687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackApproveSignDataCtx687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackApproveSignDataCtx687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackApproveSignDataCtx687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackApproveSignDataCtx687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackApproveSignDataCtx687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackApproveSignDataCtx687eda_TimerEvent(this, event); };
	QString remoteDefault() { return _remote; };
	void setRemoteDefault(QString p) { if (p != _remote) { _remote = p; remoteChanged(_remote); } };
	QString transportDefault() { return _transport; };
	void setTransportDefault(QString p) { if (p != _transport) { _transport = p; transportChanged(_transport); } };
	QString endpointDefault() { return _endpoint; };
	void setEndpointDefault(QString p) { if (p != _endpoint) { _endpoint = p; endpointChanged(_endpoint); } };
	QString fromDefault() { return _from; };
	void setFromDefault(QString p) { if (p != _from) { _from = p; fromChanged(_from); } };
	QString messageDefault() { return _message; };
	void setMessageDefault(QString p) { if (p != _message) { _message = p; messageChanged(_message); } };
	QString rawDataDefault() { return _rawData; };
	void setRawDataDefault(QString p) { if (p != _rawData) { _rawData = p; rawDataChanged(_rawData); } };
	QString hashDefault() { return _hash; };
	void setHashDefault(QString p) { if (p != _hash) { _hash = p; hashChanged(_hash); } };
	QString passwordDefault() { return _password; };
	void setPasswordDefault(QString p) { if (p != _password) { _password = p; passwordChanged(_password); } };
	QString fromSrcDefault() { return _fromSrc; };
	void setFromSrcDefault(QString p) { if (p != _fromSrc) { _fromSrc = p; fromSrcChanged(_fromSrc); } };
signals:
	void clicked(qint32 b);
	void onBack();
	void onApprove();
	void onReject();
	void edited(QString b, QString value);
	void remoteChanged(QString remote);
	void transportChanged(QString transport);
	void endpointChanged(QString endpoint);
	void fromChanged(QString from);
	void messageChanged(QString message);
	void rawDataChanged(QString rawData);
	void hashChanged(QString hash);
	void passwordChanged(QString password);
	void fromSrcChanged(QString fromSrc);
public slots:
private:
	QString _remote;
	QString _transport;
	QString _endpoint;
	QString _from;
	QString _message;
	QString _rawData;
	QString _hash;
	QString _password;
	QString _fromSrc;
};

Q_DECLARE_METATYPE(ApproveSignDataCtx687eda*)


void ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class ApproveTxCtx687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(qint32 valueUnit READ valueUnit WRITE setValueUnit NOTIFY valueUnitChanged)
Q_PROPERTY(QString remote READ remote WRITE setRemote NOTIFY remoteChanged)
Q_PROPERTY(QString transport READ transport WRITE setTransport NOTIFY transportChanged)
Q_PROPERTY(QString endpoint READ endpoint WRITE setEndpoint NOTIFY endpointChanged)
Q_PROPERTY(QString data READ data WRITE setData NOTIFY dataChanged)
Q_PROPERTY(QString from READ from WRITE setFrom NOTIFY fromChanged)
Q_PROPERTY(QString fromWarning READ fromWarning WRITE setFromWarning NOTIFY fromWarningChanged)
Q_PROPERTY(bool fromVisible READ isFromVisible WRITE setFromVisible NOTIFY fromVisibleChanged)
Q_PROPERTY(QString to READ to WRITE setTo NOTIFY toChanged)
Q_PROPERTY(QString toWarning READ toWarning WRITE setToWarning NOTIFY toWarningChanged)
Q_PROPERTY(bool toVisible READ isToVisible WRITE setToVisible NOTIFY toVisibleChanged)
Q_PROPERTY(QString gas READ gas WRITE setGas NOTIFY gasChanged)
Q_PROPERTY(QString gasPrice READ gasPrice WRITE setGasPrice NOTIFY gasPriceChanged)
Q_PROPERTY(qint32 gasPriceUnit READ gasPriceUnit WRITE setGasPriceUnit NOTIFY gasPriceUnitChanged)
Q_PROPERTY(QString nonce READ nonce WRITE setNonce NOTIFY nonceChanged)
Q_PROPERTY(QString value READ value WRITE setValue NOTIFY valueChanged)
Q_PROPERTY(QString password READ password WRITE setPassword NOTIFY passwordChanged)
Q_PROPERTY(QString fromSrc READ fromSrc WRITE setFromSrc NOTIFY fromSrcChanged)
Q_PROPERTY(QString toSrc READ toSrc WRITE setToSrc NOTIFY toSrcChanged)
Q_PROPERTY(QString diff READ diff WRITE setDiff NOTIFY diffChanged)
public:
	ApproveTxCtx687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType();ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaTypes();callbackApproveTxCtx687eda_Constructor(this);};
	void Signal_Approve() { callbackApproveTxCtx687eda_Approve(this); };
	void Signal_Reject() { callbackApproveTxCtx687eda_Reject(this); };
	void Signal_CheckTxDiff() { callbackApproveTxCtx687eda_CheckTxDiff(this); };
	void Signal_Back() { callbackApproveTxCtx687eda_Back(this); };
	void Signal_Edited(QString s, QString v) { QByteArray* ta0f149 = new QByteArray(s.toUtf8()); Moc_PackedString sPacked = { const_cast<char*>(ta0f149->prepend("WHITESPACE").constData()+10), ta0f149->size()-10, ta0f149 };QByteArray* t7a38d8 = new QByteArray(v.toUtf8()); Moc_PackedString vPacked = { const_cast<char*>(t7a38d8->prepend("WHITESPACE").constData()+10), t7a38d8->size()-10, t7a38d8 };callbackApproveTxCtx687eda_Edited(this, sPacked, vPacked); };
	void Signal_ChangeValueUnit(qint32 v) { callbackApproveTxCtx687eda_ChangeValueUnit(this, v); };
	void Signal_ChangeGasPriceUnit(qint32 v) { callbackApproveTxCtx687eda_ChangeGasPriceUnit(this, v); };
	qint32 valueUnit() { return callbackApproveTxCtx687eda_ValueUnit(this); };
	void setValueUnit(qint32 valueUnit) { callbackApproveTxCtx687eda_SetValueUnit(this, valueUnit); };
	void Signal_ValueUnitChanged(qint32 valueUnit) { callbackApproveTxCtx687eda_ValueUnitChanged(this, valueUnit); };
	QString remote() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Remote(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setRemote(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveTxCtx687eda_SetRemote(this, remotePacked); };
	void Signal_RemoteChanged(QString remote) { QByteArray* t41ffe5 = new QByteArray(remote.toUtf8()); Moc_PackedString remotePacked = { const_cast<char*>(t41ffe5->prepend("WHITESPACE").constData()+10), t41ffe5->size()-10, t41ffe5 };callbackApproveTxCtx687eda_RemoteChanged(this, remotePacked); };
	QString transport() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Transport(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setTransport(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveTxCtx687eda_SetTransport(this, transportPacked); };
	void Signal_TransportChanged(QString transport) { QByteArray* ta8e601 = new QByteArray(transport.toUtf8()); Moc_PackedString transportPacked = { const_cast<char*>(ta8e601->prepend("WHITESPACE").constData()+10), ta8e601->size()-10, ta8e601 };callbackApproveTxCtx687eda_TransportChanged(this, transportPacked); };
	QString endpoint() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Endpoint(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setEndpoint(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveTxCtx687eda_SetEndpoint(this, endpointPacked); };
	void Signal_EndpointChanged(QString endpoint) { QByteArray* te13fe4 = new QByteArray(endpoint.toUtf8()); Moc_PackedString endpointPacked = { const_cast<char*>(te13fe4->prepend("WHITESPACE").constData()+10), te13fe4->size()-10, te13fe4 };callbackApproveTxCtx687eda_EndpointChanged(this, endpointPacked); };
	QString data() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Data(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setData(QString data) { QByteArray* ta17c9a = new QByteArray(data.toUtf8()); Moc_PackedString dataPacked = { const_cast<char*>(ta17c9a->prepend("WHITESPACE").constData()+10), ta17c9a->size()-10, ta17c9a };callbackApproveTxCtx687eda_SetData(this, dataPacked); };
	void Signal_DataChanged(QString data) { QByteArray* ta17c9a = new QByteArray(data.toUtf8()); Moc_PackedString dataPacked = { const_cast<char*>(ta17c9a->prepend("WHITESPACE").constData()+10), ta17c9a->size()-10, ta17c9a };callbackApproveTxCtx687eda_DataChanged(this, dataPacked); };
	QString from() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_From(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFrom(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveTxCtx687eda_SetFrom(this, fromPacked); };
	void Signal_FromChanged(QString from) { QByteArray* t0b1e95 = new QByteArray(from.toUtf8()); Moc_PackedString fromPacked = { const_cast<char*>(t0b1e95->prepend("WHITESPACE").constData()+10), t0b1e95->size()-10, t0b1e95 };callbackApproveTxCtx687eda_FromChanged(this, fromPacked); };
	QString fromWarning() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_FromWarning(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFromWarning(QString fromWarning) { QByteArray* t388a5b = new QByteArray(fromWarning.toUtf8()); Moc_PackedString fromWarningPacked = { const_cast<char*>(t388a5b->prepend("WHITESPACE").constData()+10), t388a5b->size()-10, t388a5b };callbackApproveTxCtx687eda_SetFromWarning(this, fromWarningPacked); };
	void Signal_FromWarningChanged(QString fromWarning) { QByteArray* t388a5b = new QByteArray(fromWarning.toUtf8()); Moc_PackedString fromWarningPacked = { const_cast<char*>(t388a5b->prepend("WHITESPACE").constData()+10), t388a5b->size()-10, t388a5b };callbackApproveTxCtx687eda_FromWarningChanged(this, fromWarningPacked); };
	bool isFromVisible() { return callbackApproveTxCtx687eda_IsFromVisible(this) != 0; };
	void setFromVisible(bool fromVisible) { callbackApproveTxCtx687eda_SetFromVisible(this, fromVisible); };
	void Signal_FromVisibleChanged(bool fromVisible) { callbackApproveTxCtx687eda_FromVisibleChanged(this, fromVisible); };
	QString to() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_To(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setTo(QString to) { QByteArray* t4374aa = new QByteArray(to.toUtf8()); Moc_PackedString toPacked = { const_cast<char*>(t4374aa->prepend("WHITESPACE").constData()+10), t4374aa->size()-10, t4374aa };callbackApproveTxCtx687eda_SetTo(this, toPacked); };
	void Signal_ToChanged(QString to) { QByteArray* t4374aa = new QByteArray(to.toUtf8()); Moc_PackedString toPacked = { const_cast<char*>(t4374aa->prepend("WHITESPACE").constData()+10), t4374aa->size()-10, t4374aa };callbackApproveTxCtx687eda_ToChanged(this, toPacked); };
	QString toWarning() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_ToWarning(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setToWarning(QString toWarning) { QByteArray* t9fefd3 = new QByteArray(toWarning.toUtf8()); Moc_PackedString toWarningPacked = { const_cast<char*>(t9fefd3->prepend("WHITESPACE").constData()+10), t9fefd3->size()-10, t9fefd3 };callbackApproveTxCtx687eda_SetToWarning(this, toWarningPacked); };
	void Signal_ToWarningChanged(QString toWarning) { QByteArray* t9fefd3 = new QByteArray(toWarning.toUtf8()); Moc_PackedString toWarningPacked = { const_cast<char*>(t9fefd3->prepend("WHITESPACE").constData()+10), t9fefd3->size()-10, t9fefd3 };callbackApproveTxCtx687eda_ToWarningChanged(this, toWarningPacked); };
	bool isToVisible() { return callbackApproveTxCtx687eda_IsToVisible(this) != 0; };
	void setToVisible(bool toVisible) { callbackApproveTxCtx687eda_SetToVisible(this, toVisible); };
	void Signal_ToVisibleChanged(bool toVisible) { callbackApproveTxCtx687eda_ToVisibleChanged(this, toVisible); };
	QString gas() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Gas(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setGas(QString gas) { QByteArray* tfacafa = new QByteArray(gas.toUtf8()); Moc_PackedString gasPacked = { const_cast<char*>(tfacafa->prepend("WHITESPACE").constData()+10), tfacafa->size()-10, tfacafa };callbackApproveTxCtx687eda_SetGas(this, gasPacked); };
	void Signal_GasChanged(QString gas) { QByteArray* tfacafa = new QByteArray(gas.toUtf8()); Moc_PackedString gasPacked = { const_cast<char*>(tfacafa->prepend("WHITESPACE").constData()+10), tfacafa->size()-10, tfacafa };callbackApproveTxCtx687eda_GasChanged(this, gasPacked); };
	QString gasPrice() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_GasPrice(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setGasPrice(QString gasPrice) { QByteArray* t72824c = new QByteArray(gasPrice.toUtf8()); Moc_PackedString gasPricePacked = { const_cast<char*>(t72824c->prepend("WHITESPACE").constData()+10), t72824c->size()-10, t72824c };callbackApproveTxCtx687eda_SetGasPrice(this, gasPricePacked); };
	void Signal_GasPriceChanged(QString gasPrice) { QByteArray* t72824c = new QByteArray(gasPrice.toUtf8()); Moc_PackedString gasPricePacked = { const_cast<char*>(t72824c->prepend("WHITESPACE").constData()+10), t72824c->size()-10, t72824c };callbackApproveTxCtx687eda_GasPriceChanged(this, gasPricePacked); };
	qint32 gasPriceUnit() { return callbackApproveTxCtx687eda_GasPriceUnit(this); };
	void setGasPriceUnit(qint32 gasPriceUnit) { callbackApproveTxCtx687eda_SetGasPriceUnit(this, gasPriceUnit); };
	void Signal_GasPriceUnitChanged(qint32 gasPriceUnit) { callbackApproveTxCtx687eda_GasPriceUnitChanged(this, gasPriceUnit); };
	QString nonce() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Nonce(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setNonce(QString nonce) { QByteArray* t49afa7 = new QByteArray(nonce.toUtf8()); Moc_PackedString noncePacked = { const_cast<char*>(t49afa7->prepend("WHITESPACE").constData()+10), t49afa7->size()-10, t49afa7 };callbackApproveTxCtx687eda_SetNonce(this, noncePacked); };
	void Signal_NonceChanged(QString nonce) { QByteArray* t49afa7 = new QByteArray(nonce.toUtf8()); Moc_PackedString noncePacked = { const_cast<char*>(t49afa7->prepend("WHITESPACE").constData()+10), t49afa7->size()-10, t49afa7 };callbackApproveTxCtx687eda_NonceChanged(this, noncePacked); };
	QString value() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Value(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setValue(QString value) { QByteArray* tf32b67 = new QByteArray(value.toUtf8()); Moc_PackedString valuePacked = { const_cast<char*>(tf32b67->prepend("WHITESPACE").constData()+10), tf32b67->size()-10, tf32b67 };callbackApproveTxCtx687eda_SetValue(this, valuePacked); };
	void Signal_ValueChanged(QString value) { QByteArray* tf32b67 = new QByteArray(value.toUtf8()); Moc_PackedString valuePacked = { const_cast<char*>(tf32b67->prepend("WHITESPACE").constData()+10), tf32b67->size()-10, tf32b67 };callbackApproveTxCtx687eda_ValueChanged(this, valuePacked); };
	QString password() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Password(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setPassword(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveTxCtx687eda_SetPassword(this, passwordPacked); };
	void Signal_PasswordChanged(QString password) { QByteArray* t5baa61 = new QByteArray(password.toUtf8()); Moc_PackedString passwordPacked = { const_cast<char*>(t5baa61->prepend("WHITESPACE").constData()+10), t5baa61->size()-10, t5baa61 };callbackApproveTxCtx687eda_PasswordChanged(this, passwordPacked); };
	QString fromSrc() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_FromSrc(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setFromSrc(QString fromSrc) { QByteArray* ta8ded4 = new QByteArray(fromSrc.toUtf8()); Moc_PackedString fromSrcPacked = { const_cast<char*>(ta8ded4->prepend("WHITESPACE").constData()+10), ta8ded4->size()-10, ta8ded4 };callbackApproveTxCtx687eda_SetFromSrc(this, fromSrcPacked); };
	void Signal_FromSrcChanged(QString fromSrc) { QByteArray* ta8ded4 = new QByteArray(fromSrc.toUtf8()); Moc_PackedString fromSrcPacked = { const_cast<char*>(ta8ded4->prepend("WHITESPACE").constData()+10), ta8ded4->size()-10, ta8ded4 };callbackApproveTxCtx687eda_FromSrcChanged(this, fromSrcPacked); };
	QString toSrc() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_ToSrc(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setToSrc(QString toSrc) { QByteArray* t6f94e3 = new QByteArray(toSrc.toUtf8()); Moc_PackedString toSrcPacked = { const_cast<char*>(t6f94e3->prepend("WHITESPACE").constData()+10), t6f94e3->size()-10, t6f94e3 };callbackApproveTxCtx687eda_SetToSrc(this, toSrcPacked); };
	void Signal_ToSrcChanged(QString toSrc) { QByteArray* t6f94e3 = new QByteArray(toSrc.toUtf8()); Moc_PackedString toSrcPacked = { const_cast<char*>(t6f94e3->prepend("WHITESPACE").constData()+10), t6f94e3->size()-10, t6f94e3 };callbackApproveTxCtx687eda_ToSrcChanged(this, toSrcPacked); };
	QString diff() { return ({ Moc_PackedString tempVal = callbackApproveTxCtx687eda_Diff(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setDiff(QString diff) { QByteArray* t75a0ee = new QByteArray(diff.toUtf8()); Moc_PackedString diffPacked = { const_cast<char*>(t75a0ee->prepend("WHITESPACE").constData()+10), t75a0ee->size()-10, t75a0ee };callbackApproveTxCtx687eda_SetDiff(this, diffPacked); };
	void Signal_DiffChanged(QString diff) { QByteArray* t75a0ee = new QByteArray(diff.toUtf8()); Moc_PackedString diffPacked = { const_cast<char*>(t75a0ee->prepend("WHITESPACE").constData()+10), t75a0ee->size()-10, t75a0ee };callbackApproveTxCtx687eda_DiffChanged(this, diffPacked); };
	 ~ApproveTxCtx687eda() { callbackApproveTxCtx687eda_DestroyApproveTxCtx(this); };
	void childEvent(QChildEvent * event) { callbackApproveTxCtx687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackApproveTxCtx687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackApproveTxCtx687eda_CustomEvent(this, event); };
	void deleteLater() { callbackApproveTxCtx687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackApproveTxCtx687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackApproveTxCtx687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackApproveTxCtx687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackApproveTxCtx687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackApproveTxCtx687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackApproveTxCtx687eda_TimerEvent(this, event); };
	qint32 valueUnitDefault() { return _valueUnit; };
	void setValueUnitDefault(qint32 p) { if (p != _valueUnit) { _valueUnit = p; valueUnitChanged(_valueUnit); } };
	QString remoteDefault() { return _remote; };
	void setRemoteDefault(QString p) { if (p != _remote) { _remote = p; remoteChanged(_remote); } };
	QString transportDefault() { return _transport; };
	void setTransportDefault(QString p) { if (p != _transport) { _transport = p; transportChanged(_transport); } };
	QString endpointDefault() { return _endpoint; };
	void setEndpointDefault(QString p) { if (p != _endpoint) { _endpoint = p; endpointChanged(_endpoint); } };
	QString dataDefault() { return _data; };
	void setDataDefault(QString p) { if (p != _data) { _data = p; dataChanged(_data); } };
	QString fromDefault() { return _from; };
	void setFromDefault(QString p) { if (p != _from) { _from = p; fromChanged(_from); } };
	QString fromWarningDefault() { return _fromWarning; };
	void setFromWarningDefault(QString p) { if (p != _fromWarning) { _fromWarning = p; fromWarningChanged(_fromWarning); } };
	bool isFromVisibleDefault() { return _fromVisible; };
	void setFromVisibleDefault(bool p) { if (p != _fromVisible) { _fromVisible = p; fromVisibleChanged(_fromVisible); } };
	QString toDefault() { return _to; };
	void setToDefault(QString p) { if (p != _to) { _to = p; toChanged(_to); } };
	QString toWarningDefault() { return _toWarning; };
	void setToWarningDefault(QString p) { if (p != _toWarning) { _toWarning = p; toWarningChanged(_toWarning); } };
	bool isToVisibleDefault() { return _toVisible; };
	void setToVisibleDefault(bool p) { if (p != _toVisible) { _toVisible = p; toVisibleChanged(_toVisible); } };
	QString gasDefault() { return _gas; };
	void setGasDefault(QString p) { if (p != _gas) { _gas = p; gasChanged(_gas); } };
	QString gasPriceDefault() { return _gasPrice; };
	void setGasPriceDefault(QString p) { if (p != _gasPrice) { _gasPrice = p; gasPriceChanged(_gasPrice); } };
	qint32 gasPriceUnitDefault() { return _gasPriceUnit; };
	void setGasPriceUnitDefault(qint32 p) { if (p != _gasPriceUnit) { _gasPriceUnit = p; gasPriceUnitChanged(_gasPriceUnit); } };
	QString nonceDefault() { return _nonce; };
	void setNonceDefault(QString p) { if (p != _nonce) { _nonce = p; nonceChanged(_nonce); } };
	QString valueDefault() { return _value; };
	void setValueDefault(QString p) { if (p != _value) { _value = p; valueChanged(_value); } };
	QString passwordDefault() { return _password; };
	void setPasswordDefault(QString p) { if (p != _password) { _password = p; passwordChanged(_password); } };
	QString fromSrcDefault() { return _fromSrc; };
	void setFromSrcDefault(QString p) { if (p != _fromSrc) { _fromSrc = p; fromSrcChanged(_fromSrc); } };
	QString toSrcDefault() { return _toSrc; };
	void setToSrcDefault(QString p) { if (p != _toSrc) { _toSrc = p; toSrcChanged(_toSrc); } };
	QString diffDefault() { return _diff; };
	void setDiffDefault(QString p) { if (p != _diff) { _diff = p; diffChanged(_diff); } };
signals:
	void approve();
	void reject();
	void checkTxDiff();
	void back();
	void edited(QString s, QString v);
	void changeValueUnit(qint32 v);
	void changeGasPriceUnit(qint32 v);
	void valueUnitChanged(qint32 valueUnit);
	void remoteChanged(QString remote);
	void transportChanged(QString transport);
	void endpointChanged(QString endpoint);
	void dataChanged(QString data);
	void fromChanged(QString from);
	void fromWarningChanged(QString fromWarning);
	void fromVisibleChanged(bool fromVisible);
	void toChanged(QString to);
	void toWarningChanged(QString toWarning);
	void toVisibleChanged(bool toVisible);
	void gasChanged(QString gas);
	void gasPriceChanged(QString gasPrice);
	void gasPriceUnitChanged(qint32 gasPriceUnit);
	void nonceChanged(QString nonce);
	void valueChanged(QString value);
	void passwordChanged(QString password);
	void fromSrcChanged(QString fromSrc);
	void toSrcChanged(QString toSrc);
	void diffChanged(QString diff);
public slots:
private:
	qint32 _valueUnit;
	QString _remote;
	QString _transport;
	QString _endpoint;
	QString _data;
	QString _from;
	QString _fromWarning;
	bool _fromVisible;
	QString _to;
	QString _toWarning;
	bool _toVisible;
	QString _gas;
	QString _gasPrice;
	qint32 _gasPriceUnit;
	QString _nonce;
	QString _value;
	QString _password;
	QString _fromSrc;
	QString _toSrc;
	QString _diff;
};

Q_DECLARE_METATYPE(ApproveTxCtx687eda*)


void ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class CustomListModel687eda: public QAbstractListModel
{
Q_OBJECT
Q_PROPERTY(QString updateRequired READ updateRequired WRITE setUpdateRequired NOTIFY updateRequiredChanged)
public:
	CustomListModel687eda(QObject *parent = Q_NULLPTR) : QAbstractListModel(parent) {qRegisterMetaType<quintptr>("quintptr");CustomListModel687eda_CustomListModel687eda_QRegisterMetaType();CustomListModel687eda_CustomListModel687eda_QRegisterMetaTypes();callbackCustomListModel687eda_Constructor(this);};
	void Signal_Clear() { callbackCustomListModel687eda_Clear(this); };
	void Signal_Add(quintptr account) { callbackCustomListModel687eda_Add(this, account); };
	QString updateRequired() { return ({ Moc_PackedString tempVal = callbackCustomListModel687eda_UpdateRequired(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setUpdateRequired(QString updateRequired) { QByteArray* t4432f4 = new QByteArray(updateRequired.toUtf8()); Moc_PackedString updateRequiredPacked = { const_cast<char*>(t4432f4->prepend("WHITESPACE").constData()+10), t4432f4->size()-10, t4432f4 };callbackCustomListModel687eda_SetUpdateRequired(this, updateRequiredPacked); };
	void Signal_UpdateRequiredChanged(QString updateRequired) { QByteArray* t4432f4 = new QByteArray(updateRequired.toUtf8()); Moc_PackedString updateRequiredPacked = { const_cast<char*>(t4432f4->prepend("WHITESPACE").constData()+10), t4432f4->size()-10, t4432f4 };callbackCustomListModel687eda_UpdateRequiredChanged(this, updateRequiredPacked); };
	 ~CustomListModel687eda() { callbackCustomListModel687eda_DestroyCustomListModel(this); };
	bool dropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) { return callbackCustomListModel687eda_DropMimeData(this, const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	Qt::ItemFlags flags(const QModelIndex & index) const { return static_cast<Qt::ItemFlag>(callbackCustomListModel687eda_Flags(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	QModelIndex index(int row, int column, const QModelIndex & parent) const { return *static_cast<QModelIndex*>(callbackCustomListModel687eda_Index(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&parent))); };
	QModelIndex sibling(int row, int column, const QModelIndex & idx) const { return *static_cast<QModelIndex*>(callbackCustomListModel687eda_Sibling(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&idx))); };
	QModelIndex buddy(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackCustomListModel687eda_Buddy(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool canDropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) const { return callbackCustomListModel687eda_CanDropMimeData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	bool canFetchMore(const QModelIndex & parent) const { return callbackCustomListModel687eda_CanFetchMore(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	int columnCount(const QModelIndex & parent) const { return callbackCustomListModel687eda_ColumnCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_ColumnsAboutToBeInserted(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_ColumnsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationColumn) { callbackCustomListModel687eda_ColumnsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationColumn); };
	void Signal_ColumnsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_ColumnsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsInserted(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_ColumnsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int column) { callbackCustomListModel687eda_ColumnsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), column); };
	void Signal_ColumnsRemoved(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_ColumnsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	QVariant data(const QModelIndex & index, int role) const { return *static_cast<QVariant*>(callbackCustomListModel687eda_Data(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index), role)); };
	void Signal_DataChanged(const QModelIndex & topLeft, const QModelIndex & bottomRight, const QVector<int> & roles) { callbackCustomListModel687eda_DataChanged(this, const_cast<QModelIndex*>(&topLeft), const_cast<QModelIndex*>(&bottomRight), ({ QVector<int>* tmpValue037c88 = new QVector<int>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })); };
	void fetchMore(const QModelIndex & parent) { callbackCustomListModel687eda_FetchMore(this, const_cast<QModelIndex*>(&parent)); };
	bool hasChildren(const QModelIndex & parent) const { return callbackCustomListModel687eda_HasChildren(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	QVariant headerData(int section, Qt::Orientation orientation, int role) const { return *static_cast<QVariant*>(callbackCustomListModel687eda_HeaderData(const_cast<void*>(static_cast<const void*>(this)), section, orientation, role)); };
	void Signal_HeaderDataChanged(Qt::Orientation orientation, int first, int last) { callbackCustomListModel687eda_HeaderDataChanged(this, orientation, first, last); };
	bool insertColumns(int column, int count, const QModelIndex & parent) { return callbackCustomListModel687eda_InsertColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool insertRows(int row, int count, const QModelIndex & parent) { return callbackCustomListModel687eda_InsertRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	QMap<int, QVariant> itemData(const QModelIndex & index) const { return ({ QMap<int, QVariant>* tmpP = static_cast<QMap<int, QVariant>*>(callbackCustomListModel687eda_ItemData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); QMap<int, QVariant> tmpV = *tmpP; tmpP->~QMap(); free(tmpP); tmpV; }); };
	void Signal_LayoutAboutToBeChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackCustomListModel687eda_LayoutAboutToBeChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	void Signal_LayoutChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackCustomListModel687eda_LayoutChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	QList<QModelIndex> match(const QModelIndex & start, int role, const QVariant & value, int hits, Qt::MatchFlags flags) const { return ({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(callbackCustomListModel687eda_Match(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&start), role, const_cast<QVariant*>(&value), hits, flags)); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }); };
	QMimeData * mimeData(const QModelIndexList & indexes) const { return static_cast<QMimeData*>(callbackCustomListModel687eda_MimeData(const_cast<void*>(static_cast<const void*>(this)), ({ QList<QModelIndex>* tmpValuee0adf2 = new QList<QModelIndex>(indexes); Moc_PackedList { tmpValuee0adf2, tmpValuee0adf2->size() }; }))); };
	QStringList mimeTypes() const { return ({ Moc_PackedString tempVal = callbackCustomListModel687eda_MimeTypes(const_cast<void*>(static_cast<const void*>(this))); QStringList ret = QString::fromUtf8(tempVal.data, tempVal.len).split("¡¦!", QString::SkipEmptyParts); free(tempVal.data); ret; }); };
	void Signal_ModelAboutToBeReset() { callbackCustomListModel687eda_ModelAboutToBeReset(this); };
	void Signal_ModelReset() { callbackCustomListModel687eda_ModelReset(this); };
	bool moveColumns(const QModelIndex & sourceParent, int sourceColumn, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackCustomListModel687eda_MoveColumns(this, const_cast<QModelIndex*>(&sourceParent), sourceColumn, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	bool moveRows(const QModelIndex & sourceParent, int sourceRow, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackCustomListModel687eda_MoveRows(this, const_cast<QModelIndex*>(&sourceParent), sourceRow, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	QModelIndex parent(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackCustomListModel687eda_Parent(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool removeColumns(int column, int count, const QModelIndex & parent) { return callbackCustomListModel687eda_RemoveColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool removeRows(int row, int count, const QModelIndex & parent) { return callbackCustomListModel687eda_RemoveRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	void resetInternalData() { callbackCustomListModel687eda_ResetInternalData(this); };
	void revert() { callbackCustomListModel687eda_Revert(this); };
	QHash<int, QByteArray> roleNames() const { return ({ QHash<int, QByteArray>* tmpP = static_cast<QHash<int, QByteArray>*>(callbackCustomListModel687eda_RoleNames(const_cast<void*>(static_cast<const void*>(this)))); QHash<int, QByteArray> tmpV = *tmpP; tmpP->~QHash(); free(tmpP); tmpV; }); };
	int rowCount(const QModelIndex & parent) const { return callbackCustomListModel687eda_RowCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_RowsAboutToBeInserted(const QModelIndex & parent, int start, int end) { callbackCustomListModel687eda_RowsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), start, end); };
	void Signal_RowsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationRow) { callbackCustomListModel687eda_RowsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationRow); };
	void Signal_RowsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_RowsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsInserted(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_RowsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int row) { callbackCustomListModel687eda_RowsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), row); };
	void Signal_RowsRemoved(const QModelIndex & parent, int first, int last) { callbackCustomListModel687eda_RowsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	bool setData(const QModelIndex & index, const QVariant & value, int role) { return callbackCustomListModel687eda_SetData(this, const_cast<QModelIndex*>(&index), const_cast<QVariant*>(&value), role) != 0; };
	bool setHeaderData(int section, Qt::Orientation orientation, const QVariant & value, int role) { return callbackCustomListModel687eda_SetHeaderData(this, section, orientation, const_cast<QVariant*>(&value), role) != 0; };
	bool setItemData(const QModelIndex & index, const QMap<int, QVariant> & roles) { return callbackCustomListModel687eda_SetItemData(this, const_cast<QModelIndex*>(&index), ({ QMap<int, QVariant>* tmpValue037c88 = new QMap<int, QVariant>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })) != 0; };
	void sort(int column, Qt::SortOrder order) { callbackCustomListModel687eda_Sort(this, column, order); };
	QSize span(const QModelIndex & index) const { return *static_cast<QSize*>(callbackCustomListModel687eda_Span(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool submit() { return callbackCustomListModel687eda_Submit(this) != 0; };
	Qt::DropActions supportedDragActions() const { return static_cast<Qt::DropAction>(callbackCustomListModel687eda_SupportedDragActions(const_cast<void*>(static_cast<const void*>(this)))); };
	Qt::DropActions supportedDropActions() const { return static_cast<Qt::DropAction>(callbackCustomListModel687eda_SupportedDropActions(const_cast<void*>(static_cast<const void*>(this)))); };
	void childEvent(QChildEvent * event) { callbackCustomListModel687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackCustomListModel687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackCustomListModel687eda_CustomEvent(this, event); };
	void deleteLater() { callbackCustomListModel687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackCustomListModel687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackCustomListModel687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackCustomListModel687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackCustomListModel687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackCustomListModel687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackCustomListModel687eda_TimerEvent(this, event); };
	QString updateRequiredDefault() { return _updateRequired; };
	void setUpdateRequiredDefault(QString p) { if (p != _updateRequired) { _updateRequired = p; updateRequiredChanged(_updateRequired); } };
signals:
	void clear();
	void add(quintptr account);
	void updateRequiredChanged(QString updateRequired);
public slots:
	void callbackFromQml() { callbackCustomListModel687eda_CallbackFromQml(this); };
private:
	QString _updateRequired;
};

Q_DECLARE_METATYPE(CustomListModel687eda*)


void CustomListModel687eda_CustomListModel687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class LoginContext687eda: public QObject
{
Q_OBJECT
Q_PROPERTY(QString clefPath READ clefPath WRITE setClefPath NOTIFY clefPathChanged)
Q_PROPERTY(QString binaryHash READ binaryHash WRITE setBinaryHash NOTIFY binaryHashChanged)
Q_PROPERTY(QString error READ error WRITE setError NOTIFY errorChanged)
public:
	LoginContext687eda(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");LoginContext687eda_LoginContext687eda_QRegisterMetaType();LoginContext687eda_LoginContext687eda_QRegisterMetaTypes();callbackLoginContext687eda_Constructor(this);};
	void Signal_Start() { callbackLoginContext687eda_Start(this); };
	void Signal_CheckPath(QString b) { QByteArray* te9d71f = new QByteArray(b.toUtf8()); Moc_PackedString bPacked = { const_cast<char*>(te9d71f->prepend("WHITESPACE").constData()+10), te9d71f->size()-10, te9d71f };callbackLoginContext687eda_CheckPath(this, bPacked); };
	QString clefPath() { return ({ Moc_PackedString tempVal = callbackLoginContext687eda_ClefPath(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setClefPath(QString clefPath) { QByteArray* tea58ae = new QByteArray(clefPath.toUtf8()); Moc_PackedString clefPathPacked = { const_cast<char*>(tea58ae->prepend("WHITESPACE").constData()+10), tea58ae->size()-10, tea58ae };callbackLoginContext687eda_SetClefPath(this, clefPathPacked); };
	void Signal_ClefPathChanged(QString clefPath) { QByteArray* tea58ae = new QByteArray(clefPath.toUtf8()); Moc_PackedString clefPathPacked = { const_cast<char*>(tea58ae->prepend("WHITESPACE").constData()+10), tea58ae->size()-10, tea58ae };callbackLoginContext687eda_ClefPathChanged(this, clefPathPacked); };
	QString binaryHash() { return ({ Moc_PackedString tempVal = callbackLoginContext687eda_BinaryHash(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setBinaryHash(QString binaryHash) { QByteArray* t296c7d = new QByteArray(binaryHash.toUtf8()); Moc_PackedString binaryHashPacked = { const_cast<char*>(t296c7d->prepend("WHITESPACE").constData()+10), t296c7d->size()-10, t296c7d };callbackLoginContext687eda_SetBinaryHash(this, binaryHashPacked); };
	void Signal_BinaryHashChanged(QString binaryHash) { QByteArray* t296c7d = new QByteArray(binaryHash.toUtf8()); Moc_PackedString binaryHashPacked = { const_cast<char*>(t296c7d->prepend("WHITESPACE").constData()+10), t296c7d->size()-10, t296c7d };callbackLoginContext687eda_BinaryHashChanged(this, binaryHashPacked); };
	QString error() { return ({ Moc_PackedString tempVal = callbackLoginContext687eda_Error(this); QString ret = QString::fromUtf8(tempVal.data, tempVal.len); free(tempVal.data); ret; }); };
	void setError(QString error) { QByteArray* t11f957 = new QByteArray(error.toUtf8()); Moc_PackedString errorPacked = { const_cast<char*>(t11f957->prepend("WHITESPACE").constData()+10), t11f957->size()-10, t11f957 };callbackLoginContext687eda_SetError(this, errorPacked); };
	void Signal_ErrorChanged(QString error) { QByteArray* t11f957 = new QByteArray(error.toUtf8()); Moc_PackedString errorPacked = { const_cast<char*>(t11f957->prepend("WHITESPACE").constData()+10), t11f957->size()-10, t11f957 };callbackLoginContext687eda_ErrorChanged(this, errorPacked); };
	 ~LoginContext687eda() { callbackLoginContext687eda_DestroyLoginContext(this); };
	void childEvent(QChildEvent * event) { callbackLoginContext687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackLoginContext687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackLoginContext687eda_CustomEvent(this, event); };
	void deleteLater() { callbackLoginContext687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackLoginContext687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackLoginContext687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackLoginContext687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackLoginContext687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackLoginContext687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackLoginContext687eda_TimerEvent(this, event); };
	QString clefPathDefault() { return _clefPath; };
	void setClefPathDefault(QString p) { if (p != _clefPath) { _clefPath = p; clefPathChanged(_clefPath); } };
	QString binaryHashDefault() { return _binaryHash; };
	void setBinaryHashDefault(QString p) { if (p != _binaryHash) { _binaryHash = p; binaryHashChanged(_binaryHash); } };
	QString errorDefault() { return _error; };
	void setErrorDefault(QString p) { if (p != _error) { _error = p; errorChanged(_error); } };
signals:
	void start();
	void checkPath(QString b);
	void clefPathChanged(QString clefPath);
	void binaryHashChanged(QString binaryHash);
	void errorChanged(QString error);
public slots:
private:
	QString _clefPath;
	QString _binaryHash;
	QString _error;
};

Q_DECLARE_METATYPE(LoginContext687eda*)


void LoginContext687eda_LoginContext687eda_QRegisterMetaTypes() {
	qRegisterMetaType<QString>();
}

class TxListAccountsModel687eda: public QAbstractListModel
{
Q_OBJECT
public:
	TxListAccountsModel687eda(QObject *parent = Q_NULLPTR) : QAbstractListModel(parent) {qRegisterMetaType<quintptr>("quintptr");TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType();TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaTypes();callbackTxListAccountsModel687eda_Constructor(this);};
	void Signal_Add(QString tx) { QByteArray* t066bc1 = new QByteArray(tx.toUtf8()); Moc_PackedString txPacked = { const_cast<char*>(t066bc1->prepend("WHITESPACE").constData()+10), t066bc1->size()-10, t066bc1 };callbackTxListAccountsModel687eda_Add(this, txPacked); };
	 ~TxListAccountsModel687eda() { callbackTxListAccountsModel687eda_DestroyTxListAccountsModel(this); };
	bool dropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) { return callbackTxListAccountsModel687eda_DropMimeData(this, const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	Qt::ItemFlags flags(const QModelIndex & index) const { return static_cast<Qt::ItemFlag>(callbackTxListAccountsModel687eda_Flags(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	QModelIndex index(int row, int column, const QModelIndex & parent) const { return *static_cast<QModelIndex*>(callbackTxListAccountsModel687eda_Index(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&parent))); };
	QModelIndex sibling(int row, int column, const QModelIndex & idx) const { return *static_cast<QModelIndex*>(callbackTxListAccountsModel687eda_Sibling(const_cast<void*>(static_cast<const void*>(this)), row, column, const_cast<QModelIndex*>(&idx))); };
	QModelIndex buddy(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackTxListAccountsModel687eda_Buddy(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool canDropMimeData(const QMimeData * data, Qt::DropAction action, int row, int column, const QModelIndex & parent) const { return callbackTxListAccountsModel687eda_CanDropMimeData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QMimeData*>(data), action, row, column, const_cast<QModelIndex*>(&parent)) != 0; };
	bool canFetchMore(const QModelIndex & parent) const { return callbackTxListAccountsModel687eda_CanFetchMore(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	int columnCount(const QModelIndex & parent) const { return callbackTxListAccountsModel687eda_ColumnCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_ColumnsAboutToBeInserted(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_ColumnsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationColumn) { callbackTxListAccountsModel687eda_ColumnsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationColumn); };
	void Signal_ColumnsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_ColumnsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsInserted(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_ColumnsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_ColumnsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int column) { callbackTxListAccountsModel687eda_ColumnsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), column); };
	void Signal_ColumnsRemoved(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_ColumnsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	QVariant data(const QModelIndex & index, int role) const { return *static_cast<QVariant*>(callbackTxListAccountsModel687eda_Data(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index), role)); };
	void Signal_DataChanged(const QModelIndex & topLeft, const QModelIndex & bottomRight, const QVector<int> & roles) { callbackTxListAccountsModel687eda_DataChanged(this, const_cast<QModelIndex*>(&topLeft), const_cast<QModelIndex*>(&bottomRight), ({ QVector<int>* tmpValue037c88 = new QVector<int>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })); };
	void fetchMore(const QModelIndex & parent) { callbackTxListAccountsModel687eda_FetchMore(this, const_cast<QModelIndex*>(&parent)); };
	bool hasChildren(const QModelIndex & parent) const { return callbackTxListAccountsModel687eda_HasChildren(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)) != 0; };
	QVariant headerData(int section, Qt::Orientation orientation, int role) const { return *static_cast<QVariant*>(callbackTxListAccountsModel687eda_HeaderData(const_cast<void*>(static_cast<const void*>(this)), section, orientation, role)); };
	void Signal_HeaderDataChanged(Qt::Orientation orientation, int first, int last) { callbackTxListAccountsModel687eda_HeaderDataChanged(this, orientation, first, last); };
	bool insertColumns(int column, int count, const QModelIndex & parent) { return callbackTxListAccountsModel687eda_InsertColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool insertRows(int row, int count, const QModelIndex & parent) { return callbackTxListAccountsModel687eda_InsertRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	QMap<int, QVariant> itemData(const QModelIndex & index) const { return ({ QMap<int, QVariant>* tmpP = static_cast<QMap<int, QVariant>*>(callbackTxListAccountsModel687eda_ItemData(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); QMap<int, QVariant> tmpV = *tmpP; tmpP->~QMap(); free(tmpP); tmpV; }); };
	void Signal_LayoutAboutToBeChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackTxListAccountsModel687eda_LayoutAboutToBeChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	void Signal_LayoutChanged(const QList<QPersistentModelIndex> & parents, QAbstractItemModel::LayoutChangeHint hint) { callbackTxListAccountsModel687eda_LayoutChanged(this, ({ QList<QPersistentModelIndex>* tmpValuea664f1 = new QList<QPersistentModelIndex>(parents); Moc_PackedList { tmpValuea664f1, tmpValuea664f1->size() }; }), hint); };
	QList<QModelIndex> match(const QModelIndex & start, int role, const QVariant & value, int hits, Qt::MatchFlags flags) const { return ({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(callbackTxListAccountsModel687eda_Match(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&start), role, const_cast<QVariant*>(&value), hits, flags)); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }); };
	QMimeData * mimeData(const QModelIndexList & indexes) const { return static_cast<QMimeData*>(callbackTxListAccountsModel687eda_MimeData(const_cast<void*>(static_cast<const void*>(this)), ({ QList<QModelIndex>* tmpValuee0adf2 = new QList<QModelIndex>(indexes); Moc_PackedList { tmpValuee0adf2, tmpValuee0adf2->size() }; }))); };
	QStringList mimeTypes() const { return ({ Moc_PackedString tempVal = callbackTxListAccountsModel687eda_MimeTypes(const_cast<void*>(static_cast<const void*>(this))); QStringList ret = QString::fromUtf8(tempVal.data, tempVal.len).split("¡¦!", QString::SkipEmptyParts); free(tempVal.data); ret; }); };
	void Signal_ModelAboutToBeReset() { callbackTxListAccountsModel687eda_ModelAboutToBeReset(this); };
	void Signal_ModelReset() { callbackTxListAccountsModel687eda_ModelReset(this); };
	bool moveColumns(const QModelIndex & sourceParent, int sourceColumn, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackTxListAccountsModel687eda_MoveColumns(this, const_cast<QModelIndex*>(&sourceParent), sourceColumn, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	bool moveRows(const QModelIndex & sourceParent, int sourceRow, int count, const QModelIndex & destinationParent, int destinationChild) { return callbackTxListAccountsModel687eda_MoveRows(this, const_cast<QModelIndex*>(&sourceParent), sourceRow, count, const_cast<QModelIndex*>(&destinationParent), destinationChild) != 0; };
	QModelIndex parent(const QModelIndex & index) const { return *static_cast<QModelIndex*>(callbackTxListAccountsModel687eda_Parent(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool removeColumns(int column, int count, const QModelIndex & parent) { return callbackTxListAccountsModel687eda_RemoveColumns(this, column, count, const_cast<QModelIndex*>(&parent)) != 0; };
	bool removeRows(int row, int count, const QModelIndex & parent) { return callbackTxListAccountsModel687eda_RemoveRows(this, row, count, const_cast<QModelIndex*>(&parent)) != 0; };
	void resetInternalData() { callbackTxListAccountsModel687eda_ResetInternalData(this); };
	void revert() { callbackTxListAccountsModel687eda_Revert(this); };
	QHash<int, QByteArray> roleNames() const { return ({ QHash<int, QByteArray>* tmpP = static_cast<QHash<int, QByteArray>*>(callbackTxListAccountsModel687eda_RoleNames(const_cast<void*>(static_cast<const void*>(this)))); QHash<int, QByteArray> tmpV = *tmpP; tmpP->~QHash(); free(tmpP); tmpV; }); };
	int rowCount(const QModelIndex & parent) const { return callbackTxListAccountsModel687eda_RowCount(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&parent)); };
	void Signal_RowsAboutToBeInserted(const QModelIndex & parent, int start, int end) { callbackTxListAccountsModel687eda_RowsAboutToBeInserted(this, const_cast<QModelIndex*>(&parent), start, end); };
	void Signal_RowsAboutToBeMoved(const QModelIndex & sourceParent, int sourceStart, int sourceEnd, const QModelIndex & destinationParent, int destinationRow) { callbackTxListAccountsModel687eda_RowsAboutToBeMoved(this, const_cast<QModelIndex*>(&sourceParent), sourceStart, sourceEnd, const_cast<QModelIndex*>(&destinationParent), destinationRow); };
	void Signal_RowsAboutToBeRemoved(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_RowsAboutToBeRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsInserted(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_RowsInserted(this, const_cast<QModelIndex*>(&parent), first, last); };
	void Signal_RowsMoved(const QModelIndex & parent, int start, int end, const QModelIndex & destination, int row) { callbackTxListAccountsModel687eda_RowsMoved(this, const_cast<QModelIndex*>(&parent), start, end, const_cast<QModelIndex*>(&destination), row); };
	void Signal_RowsRemoved(const QModelIndex & parent, int first, int last) { callbackTxListAccountsModel687eda_RowsRemoved(this, const_cast<QModelIndex*>(&parent), first, last); };
	bool setData(const QModelIndex & index, const QVariant & value, int role) { return callbackTxListAccountsModel687eda_SetData(this, const_cast<QModelIndex*>(&index), const_cast<QVariant*>(&value), role) != 0; };
	bool setHeaderData(int section, Qt::Orientation orientation, const QVariant & value, int role) { return callbackTxListAccountsModel687eda_SetHeaderData(this, section, orientation, const_cast<QVariant*>(&value), role) != 0; };
	bool setItemData(const QModelIndex & index, const QMap<int, QVariant> & roles) { return callbackTxListAccountsModel687eda_SetItemData(this, const_cast<QModelIndex*>(&index), ({ QMap<int, QVariant>* tmpValue037c88 = new QMap<int, QVariant>(roles); Moc_PackedList { tmpValue037c88, tmpValue037c88->size() }; })) != 0; };
	void sort(int column, Qt::SortOrder order) { callbackTxListAccountsModel687eda_Sort(this, column, order); };
	QSize span(const QModelIndex & index) const { return *static_cast<QSize*>(callbackTxListAccountsModel687eda_Span(const_cast<void*>(static_cast<const void*>(this)), const_cast<QModelIndex*>(&index))); };
	bool submit() { return callbackTxListAccountsModel687eda_Submit(this) != 0; };
	Qt::DropActions supportedDragActions() const { return static_cast<Qt::DropAction>(callbackTxListAccountsModel687eda_SupportedDragActions(const_cast<void*>(static_cast<const void*>(this)))); };
	Qt::DropActions supportedDropActions() const { return static_cast<Qt::DropAction>(callbackTxListAccountsModel687eda_SupportedDropActions(const_cast<void*>(static_cast<const void*>(this)))); };
	void childEvent(QChildEvent * event) { callbackTxListAccountsModel687eda_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackTxListAccountsModel687eda_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackTxListAccountsModel687eda_CustomEvent(this, event); };
	void deleteLater() { callbackTxListAccountsModel687eda_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackTxListAccountsModel687eda_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackTxListAccountsModel687eda_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	bool event(QEvent * e) { return callbackTxListAccountsModel687eda_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackTxListAccountsModel687eda_EventFilter(this, watched, event) != 0; };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray* taa2c4f = new QByteArray(objectName.toUtf8()); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f->prepend("WHITESPACE").constData()+10), taa2c4f->size()-10, taa2c4f };callbackTxListAccountsModel687eda_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackTxListAccountsModel687eda_TimerEvent(this, event); };
signals:
	void add(QString tx);
public slots:
private:
};

Q_DECLARE_METATYPE(TxListAccountsModel687eda*)


void TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaTypes() {
}

void CustomListModel687eda_ConnectClear(void* ptr, long long t)
{
	QObject::connect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)()>(&CustomListModel687eda::clear), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)()>(&CustomListModel687eda::Signal_Clear), static_cast<Qt::ConnectionType>(t));
}

void CustomListModel687eda_DisconnectClear(void* ptr)
{
	QObject::disconnect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)()>(&CustomListModel687eda::clear), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)()>(&CustomListModel687eda::Signal_Clear));
}

void CustomListModel687eda_Clear(void* ptr)
{
	static_cast<CustomListModel687eda*>(ptr)->clear();
}

void CustomListModel687eda_ConnectAdd(void* ptr, long long t)
{
	QObject::connect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(quintptr)>(&CustomListModel687eda::add), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(quintptr)>(&CustomListModel687eda::Signal_Add), static_cast<Qt::ConnectionType>(t));
}

void CustomListModel687eda_DisconnectAdd(void* ptr)
{
	QObject::disconnect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(quintptr)>(&CustomListModel687eda::add), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(quintptr)>(&CustomListModel687eda::Signal_Add));
}

void CustomListModel687eda_Add(void* ptr, uintptr_t account)
{
	static_cast<CustomListModel687eda*>(ptr)->add(account);
}

void CustomListModel687eda_CallbackFromQml(void* ptr)
{
	QMetaObject::invokeMethod(static_cast<CustomListModel687eda*>(ptr), "callbackFromQml");
}

struct Moc_PackedString CustomListModel687eda_UpdateRequired(void* ptr)
{
	return ({ QByteArray* t5b7d2c = new QByteArray(static_cast<CustomListModel687eda*>(ptr)->updateRequired().toUtf8()); Moc_PackedString { const_cast<char*>(t5b7d2c->prepend("WHITESPACE").constData()+10), t5b7d2c->size()-10, t5b7d2c }; });
}

struct Moc_PackedString CustomListModel687eda_UpdateRequiredDefault(void* ptr)
{
	return ({ QByteArray* t4d3d34 = new QByteArray(static_cast<CustomListModel687eda*>(ptr)->updateRequiredDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t4d3d34->prepend("WHITESPACE").constData()+10), t4d3d34->size()-10, t4d3d34 }; });
}

void CustomListModel687eda_SetUpdateRequired(void* ptr, struct Moc_PackedString updateRequired)
{
	static_cast<CustomListModel687eda*>(ptr)->setUpdateRequired(QString::fromUtf8(updateRequired.data, updateRequired.len));
}

void CustomListModel687eda_SetUpdateRequiredDefault(void* ptr, struct Moc_PackedString updateRequired)
{
	static_cast<CustomListModel687eda*>(ptr)->setUpdateRequiredDefault(QString::fromUtf8(updateRequired.data, updateRequired.len));
}

void CustomListModel687eda_ConnectUpdateRequiredChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(QString)>(&CustomListModel687eda::updateRequiredChanged), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(QString)>(&CustomListModel687eda::Signal_UpdateRequiredChanged), static_cast<Qt::ConnectionType>(t));
}

void CustomListModel687eda_DisconnectUpdateRequiredChanged(void* ptr)
{
	QObject::disconnect(static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(QString)>(&CustomListModel687eda::updateRequiredChanged), static_cast<CustomListModel687eda*>(ptr), static_cast<void (CustomListModel687eda::*)(QString)>(&CustomListModel687eda::Signal_UpdateRequiredChanged));
}

void CustomListModel687eda_UpdateRequiredChanged(void* ptr, struct Moc_PackedString updateRequired)
{
	static_cast<CustomListModel687eda*>(ptr)->updateRequiredChanged(QString::fromUtf8(updateRequired.data, updateRequired.len));
}

int CustomListModel687eda_CustomListModel687eda_QRegisterMetaType()
{
	return qRegisterMetaType<CustomListModel687eda*>();
}

int CustomListModel687eda_CustomListModel687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<CustomListModel687eda*>(const_cast<const char*>(typeName));
}

int CustomListModel687eda_CustomListModel687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<CustomListModel687eda>();
#else
	return 0;
#endif
}

int CustomListModel687eda_CustomListModel687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<CustomListModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int CustomListModel687eda_CustomListModel687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<CustomListModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

int CustomListModel687eda_____itemData_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda_____itemData_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* CustomListModel687eda_____itemData_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int CustomListModel687eda_____roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda_____roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* CustomListModel687eda_____roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int CustomListModel687eda_____setItemData_roles_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda_____setItemData_roles_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* CustomListModel687eda_____setItemData_roles_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* CustomListModel687eda___changePersistentIndexList_from_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___changePersistentIndexList_from_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* CustomListModel687eda___changePersistentIndexList_from_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* CustomListModel687eda___changePersistentIndexList_to_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___changePersistentIndexList_to_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* CustomListModel687eda___changePersistentIndexList_to_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

int CustomListModel687eda___dataChanged_roles_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QVector<int>*>(ptr)->at(i); if (i == static_cast<QVector<int>*>(ptr)->size()-1) { static_cast<QVector<int>*>(ptr)->~QVector(); free(ptr); }; tmp; });
}

void CustomListModel687eda___dataChanged_roles_setList(void* ptr, int i)
{
	static_cast<QVector<int>*>(ptr)->append(i);
}

void* CustomListModel687eda___dataChanged_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QVector<int>();
}

void* CustomListModel687eda___itemData_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___itemData_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* CustomListModel687eda___itemData_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList CustomListModel687eda___itemData_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

void* CustomListModel687eda___layoutAboutToBeChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___layoutAboutToBeChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* CustomListModel687eda___layoutAboutToBeChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* CustomListModel687eda___layoutChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___layoutChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* CustomListModel687eda___layoutChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* CustomListModel687eda___match_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___match_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* CustomListModel687eda___match_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* CustomListModel687eda___mimeData_indexes_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___mimeData_indexes_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* CustomListModel687eda___mimeData_indexes_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* CustomListModel687eda___persistentIndexList_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___persistentIndexList_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* CustomListModel687eda___persistentIndexList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* CustomListModel687eda___roleNames_atList(void* ptr, int v, int i)
{
	return new QByteArray(({ QByteArray tmp = static_cast<QHash<int, QByteArray>*>(ptr)->value(v); if (i == static_cast<QHash<int, QByteArray>*>(ptr)->size()-1) { static_cast<QHash<int, QByteArray>*>(ptr)->~QHash(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___roleNames_setList(void* ptr, int key, void* i)
{
	static_cast<QHash<int, QByteArray>*>(ptr)->insert(key, *static_cast<QByteArray*>(i));
}

void* CustomListModel687eda___roleNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QHash<int, QByteArray>();
}

struct Moc_PackedList CustomListModel687eda___roleNames_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue7fc3bb = new QList<int>(static_cast<QHash<int, QByteArray>*>(ptr)->keys()); Moc_PackedList { tmpValue7fc3bb, tmpValue7fc3bb->size() }; });
}

void* CustomListModel687eda___setItemData_roles_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___setItemData_roles_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* CustomListModel687eda___setItemData_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList CustomListModel687eda___setItemData_roles_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

int CustomListModel687eda_____doSetRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda_____doSetRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* CustomListModel687eda_____doSetRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int CustomListModel687eda_____setRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda_____setRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* CustomListModel687eda_____setRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* CustomListModel687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* CustomListModel687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* CustomListModel687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void CustomListModel687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* CustomListModel687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* CustomListModel687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* CustomListModel687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* CustomListModel687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void CustomListModel687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* CustomListModel687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* CustomListModel687eda_NewCustomListModel(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new CustomListModel687eda(static_cast<QWindow*>(parent));
	} else {
		return new CustomListModel687eda(static_cast<QObject*>(parent));
	}
}

void CustomListModel687eda_DestroyCustomListModel(void* ptr)
{
	static_cast<CustomListModel687eda*>(ptr)->~CustomListModel687eda();
}

void CustomListModel687eda_DestroyCustomListModelDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

char CustomListModel687eda_DropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::dropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

long long CustomListModel687eda_FlagsDefault(void* ptr, void* index)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::flags(*static_cast<QModelIndex*>(index));
}

void* CustomListModel687eda_IndexDefault(void* ptr, int row, int column, void* parent)
{
	return new QModelIndex(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::index(row, column, *static_cast<QModelIndex*>(parent)));
}

void* CustomListModel687eda_SiblingDefault(void* ptr, int row, int column, void* idx)
{
	return new QModelIndex(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::sibling(row, column, *static_cast<QModelIndex*>(idx)));
}

void* CustomListModel687eda_BuddyDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::buddy(*static_cast<QModelIndex*>(index)));
}

char CustomListModel687eda_CanDropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::canDropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

char CustomListModel687eda_CanFetchMoreDefault(void* ptr, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::canFetchMore(*static_cast<QModelIndex*>(parent));
}

int CustomListModel687eda_ColumnCountDefault(void* ptr, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::columnCount(*static_cast<QModelIndex*>(parent));
}

void* CustomListModel687eda_DataDefault(void* ptr, void* index, int role)
{
	Q_UNUSED(ptr);
	Q_UNUSED(index);
	Q_UNUSED(role);

}

void CustomListModel687eda_FetchMoreDefault(void* ptr, void* parent)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::fetchMore(*static_cast<QModelIndex*>(parent));
}

char CustomListModel687eda_HasChildrenDefault(void* ptr, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::hasChildren(*static_cast<QModelIndex*>(parent));
}

void* CustomListModel687eda_HeaderDataDefault(void* ptr, int section, long long orientation, int role)
{
	return new QVariant(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::headerData(section, static_cast<Qt::Orientation>(orientation), role));
}

char CustomListModel687eda_InsertColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::insertColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char CustomListModel687eda_InsertRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::insertRows(row, count, *static_cast<QModelIndex*>(parent));
}

struct Moc_PackedList CustomListModel687eda_ItemDataDefault(void* ptr, void* index)
{
	return ({ QMap<int, QVariant>* tmpValue2395e7 = new QMap<int, QVariant>(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::itemData(*static_cast<QModelIndex*>(index))); Moc_PackedList { tmpValue2395e7, tmpValue2395e7->size() }; });
}

struct Moc_PackedList CustomListModel687eda_MatchDefault(void* ptr, void* start, int role, void* value, int hits, long long flags)
{
	return ({ QList<QModelIndex>* tmpValue47fbdc = new QList<QModelIndex>(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::match(*static_cast<QModelIndex*>(start), role, *static_cast<QVariant*>(value), hits, static_cast<Qt::MatchFlag>(flags))); Moc_PackedList { tmpValue47fbdc, tmpValue47fbdc->size() }; });
}

void* CustomListModel687eda_MimeDataDefault(void* ptr, void* indexes)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::mimeData(({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(indexes); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }));
}

struct Moc_PackedString CustomListModel687eda_MimeTypesDefault(void* ptr)
{
	return ({ QByteArray* t9c6749 = new QByteArray(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::mimeTypes().join("¡¦!").toUtf8()); Moc_PackedString { const_cast<char*>(t9c6749->prepend("WHITESPACE").constData()+10), t9c6749->size()-10, t9c6749 }; });
}

char CustomListModel687eda_MoveColumnsDefault(void* ptr, void* sourceParent, int sourceColumn, int count, void* destinationParent, int destinationChild)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::moveColumns(*static_cast<QModelIndex*>(sourceParent), sourceColumn, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

char CustomListModel687eda_MoveRowsDefault(void* ptr, void* sourceParent, int sourceRow, int count, void* destinationParent, int destinationChild)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::moveRows(*static_cast<QModelIndex*>(sourceParent), sourceRow, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

void* CustomListModel687eda_ParentDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::parent(*static_cast<QModelIndex*>(index)));
}

char CustomListModel687eda_RemoveColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::removeColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char CustomListModel687eda_RemoveRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::removeRows(row, count, *static_cast<QModelIndex*>(parent));
}

void CustomListModel687eda_ResetInternalDataDefault(void* ptr)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::resetInternalData();
}

void CustomListModel687eda_RevertDefault(void* ptr)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::revert();
}

struct Moc_PackedList CustomListModel687eda_RoleNamesDefault(void* ptr)
{
	return ({ QHash<int, QByteArray>* tmpValue08990a = new QHash<int, QByteArray>(static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::roleNames()); Moc_PackedList { tmpValue08990a, tmpValue08990a->size() }; });
}

int CustomListModel687eda_RowCountDefault(void* ptr, void* parent)
{
	Q_UNUSED(ptr);
	Q_UNUSED(parent);

}

char CustomListModel687eda_SetDataDefault(void* ptr, void* index, void* value, int role)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::setData(*static_cast<QModelIndex*>(index), *static_cast<QVariant*>(value), role);
}

char CustomListModel687eda_SetHeaderDataDefault(void* ptr, int section, long long orientation, void* value, int role)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::setHeaderData(section, static_cast<Qt::Orientation>(orientation), *static_cast<QVariant*>(value), role);
}

char CustomListModel687eda_SetItemDataDefault(void* ptr, void* index, void* roles)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::setItemData(*static_cast<QModelIndex*>(index), *static_cast<QMap<int, QVariant>*>(roles));
}

void CustomListModel687eda_SortDefault(void* ptr, int column, long long order)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::sort(column, static_cast<Qt::SortOrder>(order));
}

void* CustomListModel687eda_SpanDefault(void* ptr, void* index)
{
	return ({ QSize tmpValue = static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::span(*static_cast<QModelIndex*>(index)); new QSize(tmpValue.width(), tmpValue.height()); });
}

char CustomListModel687eda_SubmitDefault(void* ptr)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::submit();
}

long long CustomListModel687eda_SupportedDragActionsDefault(void* ptr)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::supportedDragActions();
}

long long CustomListModel687eda_SupportedDropActionsDefault(void* ptr)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::supportedDropActions();
}

void CustomListModel687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::childEvent(static_cast<QChildEvent*>(event));
}

void CustomListModel687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void CustomListModel687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::customEvent(static_cast<QEvent*>(event));
}

void CustomListModel687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::deleteLater();
}

void CustomListModel687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char CustomListModel687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::event(static_cast<QEvent*>(e));
}

char CustomListModel687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void CustomListModel687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<CustomListModel687eda*>(ptr)->QAbstractListModel::timerEvent(static_cast<QTimerEvent*>(event));
}

void LoginContext687eda_ConnectStart(void* ptr, long long t)
{
	QObject::connect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)()>(&LoginContext687eda::start), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)()>(&LoginContext687eda::Signal_Start), static_cast<Qt::ConnectionType>(t));
}

void LoginContext687eda_DisconnectStart(void* ptr)
{
	QObject::disconnect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)()>(&LoginContext687eda::start), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)()>(&LoginContext687eda::Signal_Start));
}

void LoginContext687eda_Start(void* ptr)
{
	static_cast<LoginContext687eda*>(ptr)->start();
}

void LoginContext687eda_ConnectCheckPath(void* ptr, long long t)
{
	QObject::connect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::checkPath), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_CheckPath), static_cast<Qt::ConnectionType>(t));
}

void LoginContext687eda_DisconnectCheckPath(void* ptr)
{
	QObject::disconnect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::checkPath), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_CheckPath));
}

void LoginContext687eda_CheckPath(void* ptr, struct Moc_PackedString b)
{
	static_cast<LoginContext687eda*>(ptr)->checkPath(QString::fromUtf8(b.data, b.len));
}

struct Moc_PackedString LoginContext687eda_ClefPath(void* ptr)
{
	return ({ QByteArray* t5ef319 = new QByteArray(static_cast<LoginContext687eda*>(ptr)->clefPath().toUtf8()); Moc_PackedString { const_cast<char*>(t5ef319->prepend("WHITESPACE").constData()+10), t5ef319->size()-10, t5ef319 }; });
}

struct Moc_PackedString LoginContext687eda_ClefPathDefault(void* ptr)
{
	return ({ QByteArray* tdf5990 = new QByteArray(static_cast<LoginContext687eda*>(ptr)->clefPathDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tdf5990->prepend("WHITESPACE").constData()+10), tdf5990->size()-10, tdf5990 }; });
}

void LoginContext687eda_SetClefPath(void* ptr, struct Moc_PackedString clefPath)
{
	static_cast<LoginContext687eda*>(ptr)->setClefPath(QString::fromUtf8(clefPath.data, clefPath.len));
}

void LoginContext687eda_SetClefPathDefault(void* ptr, struct Moc_PackedString clefPath)
{
	static_cast<LoginContext687eda*>(ptr)->setClefPathDefault(QString::fromUtf8(clefPath.data, clefPath.len));
}

void LoginContext687eda_ConnectClefPathChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::clefPathChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_ClefPathChanged), static_cast<Qt::ConnectionType>(t));
}

void LoginContext687eda_DisconnectClefPathChanged(void* ptr)
{
	QObject::disconnect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::clefPathChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_ClefPathChanged));
}

void LoginContext687eda_ClefPathChanged(void* ptr, struct Moc_PackedString clefPath)
{
	static_cast<LoginContext687eda*>(ptr)->clefPathChanged(QString::fromUtf8(clefPath.data, clefPath.len));
}

struct Moc_PackedString LoginContext687eda_BinaryHash(void* ptr)
{
	return ({ QByteArray* t358676 = new QByteArray(static_cast<LoginContext687eda*>(ptr)->binaryHash().toUtf8()); Moc_PackedString { const_cast<char*>(t358676->prepend("WHITESPACE").constData()+10), t358676->size()-10, t358676 }; });
}

struct Moc_PackedString LoginContext687eda_BinaryHashDefault(void* ptr)
{
	return ({ QByteArray* tab0c1d = new QByteArray(static_cast<LoginContext687eda*>(ptr)->binaryHashDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tab0c1d->prepend("WHITESPACE").constData()+10), tab0c1d->size()-10, tab0c1d }; });
}

void LoginContext687eda_SetBinaryHash(void* ptr, struct Moc_PackedString binaryHash)
{
	static_cast<LoginContext687eda*>(ptr)->setBinaryHash(QString::fromUtf8(binaryHash.data, binaryHash.len));
}

void LoginContext687eda_SetBinaryHashDefault(void* ptr, struct Moc_PackedString binaryHash)
{
	static_cast<LoginContext687eda*>(ptr)->setBinaryHashDefault(QString::fromUtf8(binaryHash.data, binaryHash.len));
}

void LoginContext687eda_ConnectBinaryHashChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::binaryHashChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_BinaryHashChanged), static_cast<Qt::ConnectionType>(t));
}

void LoginContext687eda_DisconnectBinaryHashChanged(void* ptr)
{
	QObject::disconnect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::binaryHashChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_BinaryHashChanged));
}

void LoginContext687eda_BinaryHashChanged(void* ptr, struct Moc_PackedString binaryHash)
{
	static_cast<LoginContext687eda*>(ptr)->binaryHashChanged(QString::fromUtf8(binaryHash.data, binaryHash.len));
}

struct Moc_PackedString LoginContext687eda_Error(void* ptr)
{
	return ({ QByteArray* t8c43ca = new QByteArray(static_cast<LoginContext687eda*>(ptr)->error().toUtf8()); Moc_PackedString { const_cast<char*>(t8c43ca->prepend("WHITESPACE").constData()+10), t8c43ca->size()-10, t8c43ca }; });
}

struct Moc_PackedString LoginContext687eda_ErrorDefault(void* ptr)
{
	return ({ QByteArray* ta84fb8 = new QByteArray(static_cast<LoginContext687eda*>(ptr)->errorDefault().toUtf8()); Moc_PackedString { const_cast<char*>(ta84fb8->prepend("WHITESPACE").constData()+10), ta84fb8->size()-10, ta84fb8 }; });
}

void LoginContext687eda_SetError(void* ptr, struct Moc_PackedString error)
{
	static_cast<LoginContext687eda*>(ptr)->setError(QString::fromUtf8(error.data, error.len));
}

void LoginContext687eda_SetErrorDefault(void* ptr, struct Moc_PackedString error)
{
	static_cast<LoginContext687eda*>(ptr)->setErrorDefault(QString::fromUtf8(error.data, error.len));
}

void LoginContext687eda_ConnectErrorChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::errorChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_ErrorChanged), static_cast<Qt::ConnectionType>(t));
}

void LoginContext687eda_DisconnectErrorChanged(void* ptr)
{
	QObject::disconnect(static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::errorChanged), static_cast<LoginContext687eda*>(ptr), static_cast<void (LoginContext687eda::*)(QString)>(&LoginContext687eda::Signal_ErrorChanged));
}

void LoginContext687eda_ErrorChanged(void* ptr, struct Moc_PackedString error)
{
	static_cast<LoginContext687eda*>(ptr)->errorChanged(QString::fromUtf8(error.data, error.len));
}

int LoginContext687eda_LoginContext687eda_QRegisterMetaType()
{
	return qRegisterMetaType<LoginContext687eda*>();
}

int LoginContext687eda_LoginContext687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<LoginContext687eda*>(const_cast<const char*>(typeName));
}

int LoginContext687eda_LoginContext687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<LoginContext687eda>();
#else
	return 0;
#endif
}

int LoginContext687eda_LoginContext687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<LoginContext687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int LoginContext687eda_LoginContext687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<LoginContext687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* LoginContext687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void LoginContext687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* LoginContext687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* LoginContext687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void LoginContext687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* LoginContext687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* LoginContext687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void LoginContext687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* LoginContext687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* LoginContext687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void LoginContext687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* LoginContext687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* LoginContext687eda_NewLoginContext(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new LoginContext687eda(static_cast<QWindow*>(parent));
	} else {
		return new LoginContext687eda(static_cast<QObject*>(parent));
	}
}

void LoginContext687eda_DestroyLoginContext(void* ptr)
{
	static_cast<LoginContext687eda*>(ptr)->~LoginContext687eda();
}

void LoginContext687eda_DestroyLoginContextDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void LoginContext687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void LoginContext687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void LoginContext687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void LoginContext687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::deleteLater();
}

void LoginContext687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char LoginContext687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<LoginContext687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char LoginContext687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<LoginContext687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void LoginContext687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<LoginContext687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void TxListModel687eda_ConnectClear(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)()>(&TxListModel687eda::clear), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)()>(&TxListModel687eda::Signal_Clear), static_cast<Qt::ConnectionType>(t));
}

void TxListModel687eda_DisconnectClear(void* ptr)
{
	QObject::disconnect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)()>(&TxListModel687eda::clear), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)()>(&TxListModel687eda::Signal_Clear));
}

void TxListModel687eda_Clear(void* ptr)
{
	static_cast<TxListModel687eda*>(ptr)->clear();
}

void TxListModel687eda_ConnectAdd(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(quintptr)>(&TxListModel687eda::add), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(quintptr)>(&TxListModel687eda::Signal_Add), static_cast<Qt::ConnectionType>(t));
}

void TxListModel687eda_DisconnectAdd(void* ptr)
{
	QObject::disconnect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(quintptr)>(&TxListModel687eda::add), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(quintptr)>(&TxListModel687eda::Signal_Add));
}

void TxListModel687eda_Add(void* ptr, uintptr_t tx)
{
	static_cast<TxListModel687eda*>(ptr)->add(tx);
}

void TxListModel687eda_ConnectRemove(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(qint32)>(&TxListModel687eda::remove), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(qint32)>(&TxListModel687eda::Signal_Remove), static_cast<Qt::ConnectionType>(t));
}

void TxListModel687eda_DisconnectRemove(void* ptr)
{
	QObject::disconnect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(qint32)>(&TxListModel687eda::remove), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(qint32)>(&TxListModel687eda::Signal_Remove));
}

void TxListModel687eda_Remove(void* ptr, int i)
{
	static_cast<TxListModel687eda*>(ptr)->remove(i);
}

char TxListModel687eda_IsEmpty(void* ptr)
{
	return static_cast<TxListModel687eda*>(ptr)->isEmpty();
}

char TxListModel687eda_IsEmptyDefault(void* ptr)
{
	return static_cast<TxListModel687eda*>(ptr)->isEmptyDefault();
}

void TxListModel687eda_SetIsEmpty(void* ptr, char isEmpty)
{
	static_cast<TxListModel687eda*>(ptr)->setIsEmpty(isEmpty != 0);
}

void TxListModel687eda_SetIsEmptyDefault(void* ptr, char isEmpty)
{
	static_cast<TxListModel687eda*>(ptr)->setIsEmptyDefault(isEmpty != 0);
}

void TxListModel687eda_ConnectIsEmptyChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(bool)>(&TxListModel687eda::isEmptyChanged), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(bool)>(&TxListModel687eda::Signal_IsEmptyChanged), static_cast<Qt::ConnectionType>(t));
}

void TxListModel687eda_DisconnectIsEmptyChanged(void* ptr)
{
	QObject::disconnect(static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(bool)>(&TxListModel687eda::isEmptyChanged), static_cast<TxListModel687eda*>(ptr), static_cast<void (TxListModel687eda::*)(bool)>(&TxListModel687eda::Signal_IsEmptyChanged));
}

void TxListModel687eda_IsEmptyChanged(void* ptr, char isEmpty)
{
	static_cast<TxListModel687eda*>(ptr)->isEmptyChanged(isEmpty != 0);
}

int TxListModel687eda_TxListModel687eda_QRegisterMetaType()
{
	return qRegisterMetaType<TxListModel687eda*>();
}

int TxListModel687eda_TxListModel687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<TxListModel687eda*>(const_cast<const char*>(typeName));
}

int TxListModel687eda_TxListModel687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListModel687eda>();
#else
	return 0;
#endif
}

int TxListModel687eda_TxListModel687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int TxListModel687eda_TxListModel687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<TxListModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

int TxListModel687eda_____itemData_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda_____itemData_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListModel687eda_____itemData_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListModel687eda_____roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda_____roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListModel687eda_____roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListModel687eda_____setItemData_roles_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda_____setItemData_roles_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListModel687eda_____setItemData_roles_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* TxListModel687eda___changePersistentIndexList_from_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___changePersistentIndexList_from_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListModel687eda___changePersistentIndexList_from_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListModel687eda___changePersistentIndexList_to_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___changePersistentIndexList_to_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListModel687eda___changePersistentIndexList_to_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

int TxListModel687eda___dataChanged_roles_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QVector<int>*>(ptr)->at(i); if (i == static_cast<QVector<int>*>(ptr)->size()-1) { static_cast<QVector<int>*>(ptr)->~QVector(); free(ptr); }; tmp; });
}

void TxListModel687eda___dataChanged_roles_setList(void* ptr, int i)
{
	static_cast<QVector<int>*>(ptr)->append(i);
}

void* TxListModel687eda___dataChanged_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QVector<int>();
}

void* TxListModel687eda___itemData_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void TxListModel687eda___itemData_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* TxListModel687eda___itemData_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList TxListModel687eda___itemData_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

void* TxListModel687eda___layoutAboutToBeChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___layoutAboutToBeChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* TxListModel687eda___layoutAboutToBeChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* TxListModel687eda___layoutChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___layoutChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* TxListModel687eda___layoutChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* TxListModel687eda___match_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___match_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListModel687eda___match_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListModel687eda___mimeData_indexes_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___mimeData_indexes_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListModel687eda___mimeData_indexes_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListModel687eda___persistentIndexList_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___persistentIndexList_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListModel687eda___persistentIndexList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListModel687eda___roleNames_atList(void* ptr, int v, int i)
{
	return new QByteArray(({ QByteArray tmp = static_cast<QHash<int, QByteArray>*>(ptr)->value(v); if (i == static_cast<QHash<int, QByteArray>*>(ptr)->size()-1) { static_cast<QHash<int, QByteArray>*>(ptr)->~QHash(); free(ptr); }; tmp; }));
}

void TxListModel687eda___roleNames_setList(void* ptr, int key, void* i)
{
	static_cast<QHash<int, QByteArray>*>(ptr)->insert(key, *static_cast<QByteArray*>(i));
}

void* TxListModel687eda___roleNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QHash<int, QByteArray>();
}

struct Moc_PackedList TxListModel687eda___roleNames_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue7fc3bb = new QList<int>(static_cast<QHash<int, QByteArray>*>(ptr)->keys()); Moc_PackedList { tmpValue7fc3bb, tmpValue7fc3bb->size() }; });
}

void* TxListModel687eda___setItemData_roles_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void TxListModel687eda___setItemData_roles_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* TxListModel687eda___setItemData_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList TxListModel687eda___setItemData_roles_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

int TxListModel687eda_____doSetRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda_____doSetRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListModel687eda_____doSetRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListModel687eda_____setRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda_____setRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListModel687eda_____setRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* TxListModel687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListModel687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* TxListModel687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListModel687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* TxListModel687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* TxListModel687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListModel687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListModel687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListModel687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListModel687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListModel687eda_NewTxListModel(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new TxListModel687eda(static_cast<QWindow*>(parent));
	} else {
		return new TxListModel687eda(static_cast<QObject*>(parent));
	}
}

void TxListModel687eda_DestroyTxListModel(void* ptr)
{
	static_cast<TxListModel687eda*>(ptr)->~TxListModel687eda();
}

void TxListModel687eda_DestroyTxListModelDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

char TxListModel687eda_DropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::dropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

long long TxListModel687eda_FlagsDefault(void* ptr, void* index)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::flags(*static_cast<QModelIndex*>(index));
}

void* TxListModel687eda_IndexDefault(void* ptr, int row, int column, void* parent)
{
	return new QModelIndex(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::index(row, column, *static_cast<QModelIndex*>(parent)));
}

void* TxListModel687eda_SiblingDefault(void* ptr, int row, int column, void* idx)
{
	return new QModelIndex(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::sibling(row, column, *static_cast<QModelIndex*>(idx)));
}

void* TxListModel687eda_BuddyDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::buddy(*static_cast<QModelIndex*>(index)));
}

char TxListModel687eda_CanDropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::canDropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

char TxListModel687eda_CanFetchMoreDefault(void* ptr, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::canFetchMore(*static_cast<QModelIndex*>(parent));
}

int TxListModel687eda_ColumnCountDefault(void* ptr, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::columnCount(*static_cast<QModelIndex*>(parent));
}

void* TxListModel687eda_DataDefault(void* ptr, void* index, int role)
{
	Q_UNUSED(ptr);
	Q_UNUSED(index);
	Q_UNUSED(role);

}

void TxListModel687eda_FetchMoreDefault(void* ptr, void* parent)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::fetchMore(*static_cast<QModelIndex*>(parent));
}

char TxListModel687eda_HasChildrenDefault(void* ptr, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::hasChildren(*static_cast<QModelIndex*>(parent));
}

void* TxListModel687eda_HeaderDataDefault(void* ptr, int section, long long orientation, int role)
{
	return new QVariant(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::headerData(section, static_cast<Qt::Orientation>(orientation), role));
}

char TxListModel687eda_InsertColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::insertColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char TxListModel687eda_InsertRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::insertRows(row, count, *static_cast<QModelIndex*>(parent));
}

struct Moc_PackedList TxListModel687eda_ItemDataDefault(void* ptr, void* index)
{
	return ({ QMap<int, QVariant>* tmpValued07f70 = new QMap<int, QVariant>(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::itemData(*static_cast<QModelIndex*>(index))); Moc_PackedList { tmpValued07f70, tmpValued07f70->size() }; });
}

struct Moc_PackedList TxListModel687eda_MatchDefault(void* ptr, void* start, int role, void* value, int hits, long long flags)
{
	return ({ QList<QModelIndex>* tmpValue103768 = new QList<QModelIndex>(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::match(*static_cast<QModelIndex*>(start), role, *static_cast<QVariant*>(value), hits, static_cast<Qt::MatchFlag>(flags))); Moc_PackedList { tmpValue103768, tmpValue103768->size() }; });
}

void* TxListModel687eda_MimeDataDefault(void* ptr, void* indexes)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::mimeData(({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(indexes); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }));
}

struct Moc_PackedString TxListModel687eda_MimeTypesDefault(void* ptr)
{
	return ({ QByteArray* t525b22 = new QByteArray(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::mimeTypes().join("¡¦!").toUtf8()); Moc_PackedString { const_cast<char*>(t525b22->prepend("WHITESPACE").constData()+10), t525b22->size()-10, t525b22 }; });
}

char TxListModel687eda_MoveColumnsDefault(void* ptr, void* sourceParent, int sourceColumn, int count, void* destinationParent, int destinationChild)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::moveColumns(*static_cast<QModelIndex*>(sourceParent), sourceColumn, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

char TxListModel687eda_MoveRowsDefault(void* ptr, void* sourceParent, int sourceRow, int count, void* destinationParent, int destinationChild)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::moveRows(*static_cast<QModelIndex*>(sourceParent), sourceRow, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

void* TxListModel687eda_ParentDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::parent(*static_cast<QModelIndex*>(index)));
}

char TxListModel687eda_RemoveColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::removeColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char TxListModel687eda_RemoveRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::removeRows(row, count, *static_cast<QModelIndex*>(parent));
}

void TxListModel687eda_ResetInternalDataDefault(void* ptr)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::resetInternalData();
}

void TxListModel687eda_RevertDefault(void* ptr)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::revert();
}

struct Moc_PackedList TxListModel687eda_RoleNamesDefault(void* ptr)
{
	return ({ QHash<int, QByteArray>* tmpValueedb100 = new QHash<int, QByteArray>(static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::roleNames()); Moc_PackedList { tmpValueedb100, tmpValueedb100->size() }; });
}

int TxListModel687eda_RowCountDefault(void* ptr, void* parent)
{
	Q_UNUSED(ptr);
	Q_UNUSED(parent);

}

char TxListModel687eda_SetDataDefault(void* ptr, void* index, void* value, int role)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::setData(*static_cast<QModelIndex*>(index), *static_cast<QVariant*>(value), role);
}

char TxListModel687eda_SetHeaderDataDefault(void* ptr, int section, long long orientation, void* value, int role)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::setHeaderData(section, static_cast<Qt::Orientation>(orientation), *static_cast<QVariant*>(value), role);
}

char TxListModel687eda_SetItemDataDefault(void* ptr, void* index, void* roles)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::setItemData(*static_cast<QModelIndex*>(index), *static_cast<QMap<int, QVariant>*>(roles));
}

void TxListModel687eda_SortDefault(void* ptr, int column, long long order)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::sort(column, static_cast<Qt::SortOrder>(order));
}

void* TxListModel687eda_SpanDefault(void* ptr, void* index)
{
	return ({ QSize tmpValue = static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::span(*static_cast<QModelIndex*>(index)); new QSize(tmpValue.width(), tmpValue.height()); });
}

char TxListModel687eda_SubmitDefault(void* ptr)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::submit();
}

long long TxListModel687eda_SupportedDragActionsDefault(void* ptr)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::supportedDragActions();
}

long long TxListModel687eda_SupportedDropActionsDefault(void* ptr)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::supportedDropActions();
}

void TxListModel687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::childEvent(static_cast<QChildEvent*>(event));
}

void TxListModel687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void TxListModel687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::customEvent(static_cast<QEvent*>(event));
}

void TxListModel687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::deleteLater();
}

void TxListModel687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char TxListModel687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::event(static_cast<QEvent*>(e));
}

char TxListModel687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void TxListModel687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<TxListModel687eda*>(ptr)->QAbstractListModel::timerEvent(static_cast<QTimerEvent*>(event));
}

void ApproveSignDataCtx687eda_ConnectClicked(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(qint32)>(&ApproveSignDataCtx687eda::clicked), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(qint32)>(&ApproveSignDataCtx687eda::Signal_Clicked), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectClicked(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(qint32)>(&ApproveSignDataCtx687eda::clicked), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(qint32)>(&ApproveSignDataCtx687eda::Signal_Clicked));
}

void ApproveSignDataCtx687eda_Clicked(void* ptr, int b)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->clicked(b);
}

void ApproveSignDataCtx687eda_ConnectOnBack(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onBack), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnBack), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectOnBack(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onBack), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnBack));
}

void ApproveSignDataCtx687eda_OnBack(void* ptr)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->onBack();
}

void ApproveSignDataCtx687eda_ConnectOnApprove(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onApprove), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnApprove), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectOnApprove(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onApprove), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnApprove));
}

void ApproveSignDataCtx687eda_OnApprove(void* ptr)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->onApprove();
}

void ApproveSignDataCtx687eda_ConnectOnReject(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onReject), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnReject), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectOnReject(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::onReject), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)()>(&ApproveSignDataCtx687eda::Signal_OnReject));
}

void ApproveSignDataCtx687eda_OnReject(void* ptr)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->onReject();
}

void ApproveSignDataCtx687eda_ConnectEdited(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString, QString)>(&ApproveSignDataCtx687eda::edited), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString, QString)>(&ApproveSignDataCtx687eda::Signal_Edited), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectEdited(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString, QString)>(&ApproveSignDataCtx687eda::edited), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString, QString)>(&ApproveSignDataCtx687eda::Signal_Edited));
}

void ApproveSignDataCtx687eda_Edited(void* ptr, struct Moc_PackedString b, struct Moc_PackedString value)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->edited(QString::fromUtf8(b.data, b.len), QString::fromUtf8(value.data, value.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Remote(void* ptr)
{
	return ({ QByteArray* td166de = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->remote().toUtf8()); Moc_PackedString { const_cast<char*>(td166de->prepend("WHITESPACE").constData()+10), td166de->size()-10, td166de }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_RemoteDefault(void* ptr)
{
	return ({ QByteArray* t110aa3 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->remoteDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t110aa3->prepend("WHITESPACE").constData()+10), t110aa3->size()-10, t110aa3 }; });
}

void ApproveSignDataCtx687eda_SetRemote(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setRemote(QString::fromUtf8(remote.data, remote.len));
}

void ApproveSignDataCtx687eda_SetRemoteDefault(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setRemoteDefault(QString::fromUtf8(remote.data, remote.len));
}

void ApproveSignDataCtx687eda_ConnectRemoteChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::remoteChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_RemoteChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectRemoteChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::remoteChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_RemoteChanged));
}

void ApproveSignDataCtx687eda_RemoteChanged(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->remoteChanged(QString::fromUtf8(remote.data, remote.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Transport(void* ptr)
{
	return ({ QByteArray* tf8900c = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->transport().toUtf8()); Moc_PackedString { const_cast<char*>(tf8900c->prepend("WHITESPACE").constData()+10), tf8900c->size()-10, tf8900c }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_TransportDefault(void* ptr)
{
	return ({ QByteArray* t04ad21 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->transportDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t04ad21->prepend("WHITESPACE").constData()+10), t04ad21->size()-10, t04ad21 }; });
}

void ApproveSignDataCtx687eda_SetTransport(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setTransport(QString::fromUtf8(transport.data, transport.len));
}

void ApproveSignDataCtx687eda_SetTransportDefault(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setTransportDefault(QString::fromUtf8(transport.data, transport.len));
}

void ApproveSignDataCtx687eda_ConnectTransportChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::transportChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_TransportChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectTransportChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::transportChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_TransportChanged));
}

void ApproveSignDataCtx687eda_TransportChanged(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->transportChanged(QString::fromUtf8(transport.data, transport.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Endpoint(void* ptr)
{
	return ({ QByteArray* tb0291e = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->endpoint().toUtf8()); Moc_PackedString { const_cast<char*>(tb0291e->prepend("WHITESPACE").constData()+10), tb0291e->size()-10, tb0291e }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_EndpointDefault(void* ptr)
{
	return ({ QByteArray* t55e053 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->endpointDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t55e053->prepend("WHITESPACE").constData()+10), t55e053->size()-10, t55e053 }; });
}

void ApproveSignDataCtx687eda_SetEndpoint(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setEndpoint(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveSignDataCtx687eda_SetEndpointDefault(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setEndpointDefault(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveSignDataCtx687eda_ConnectEndpointChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::endpointChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_EndpointChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectEndpointChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::endpointChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_EndpointChanged));
}

void ApproveSignDataCtx687eda_EndpointChanged(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->endpointChanged(QString::fromUtf8(endpoint.data, endpoint.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_From(void* ptr)
{
	return ({ QByteArray* t1ec6d2 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->from().toUtf8()); Moc_PackedString { const_cast<char*>(t1ec6d2->prepend("WHITESPACE").constData()+10), t1ec6d2->size()-10, t1ec6d2 }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_FromDefault(void* ptr)
{
	return ({ QByteArray* t8143a5 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->fromDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t8143a5->prepend("WHITESPACE").constData()+10), t8143a5->size()-10, t8143a5 }; });
}

void ApproveSignDataCtx687eda_SetFrom(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setFrom(QString::fromUtf8(from.data, from.len));
}

void ApproveSignDataCtx687eda_SetFromDefault(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setFromDefault(QString::fromUtf8(from.data, from.len));
}

void ApproveSignDataCtx687eda_ConnectFromChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::fromChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_FromChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectFromChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::fromChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_FromChanged));
}

void ApproveSignDataCtx687eda_FromChanged(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->fromChanged(QString::fromUtf8(from.data, from.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Message(void* ptr)
{
	return ({ QByteArray* tf8c345 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->message().toUtf8()); Moc_PackedString { const_cast<char*>(tf8c345->prepend("WHITESPACE").constData()+10), tf8c345->size()-10, tf8c345 }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_MessageDefault(void* ptr)
{
	return ({ QByteArray* t59a7b5 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->messageDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t59a7b5->prepend("WHITESPACE").constData()+10), t59a7b5->size()-10, t59a7b5 }; });
}

void ApproveSignDataCtx687eda_SetMessage(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setMessage(QString::fromUtf8(message.data, message.len));
}

void ApproveSignDataCtx687eda_SetMessageDefault(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setMessageDefault(QString::fromUtf8(message.data, message.len));
}

void ApproveSignDataCtx687eda_ConnectMessageChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::messageChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_MessageChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectMessageChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::messageChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_MessageChanged));
}

void ApproveSignDataCtx687eda_MessageChanged(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->messageChanged(QString::fromUtf8(message.data, message.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_RawData(void* ptr)
{
	return ({ QByteArray* t71547a = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->rawData().toUtf8()); Moc_PackedString { const_cast<char*>(t71547a->prepend("WHITESPACE").constData()+10), t71547a->size()-10, t71547a }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_RawDataDefault(void* ptr)
{
	return ({ QByteArray* tf54367 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->rawDataDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tf54367->prepend("WHITESPACE").constData()+10), tf54367->size()-10, tf54367 }; });
}

void ApproveSignDataCtx687eda_SetRawData(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setRawData(QString::fromUtf8(rawData.data, rawData.len));
}

void ApproveSignDataCtx687eda_SetRawDataDefault(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setRawDataDefault(QString::fromUtf8(rawData.data, rawData.len));
}

void ApproveSignDataCtx687eda_ConnectRawDataChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::rawDataChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_RawDataChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectRawDataChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::rawDataChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_RawDataChanged));
}

void ApproveSignDataCtx687eda_RawDataChanged(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->rawDataChanged(QString::fromUtf8(rawData.data, rawData.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Hash(void* ptr)
{
	return ({ QByteArray* tf9f901 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->hash().toUtf8()); Moc_PackedString { const_cast<char*>(tf9f901->prepend("WHITESPACE").constData()+10), tf9f901->size()-10, tf9f901 }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_HashDefault(void* ptr)
{
	return ({ QByteArray* tcc83f5 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->hashDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tcc83f5->prepend("WHITESPACE").constData()+10), tcc83f5->size()-10, tcc83f5 }; });
}

void ApproveSignDataCtx687eda_SetHash(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setHash(QString::fromUtf8(hash.data, hash.len));
}

void ApproveSignDataCtx687eda_SetHashDefault(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setHashDefault(QString::fromUtf8(hash.data, hash.len));
}

void ApproveSignDataCtx687eda_ConnectHashChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::hashChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_HashChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectHashChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::hashChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_HashChanged));
}

void ApproveSignDataCtx687eda_HashChanged(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->hashChanged(QString::fromUtf8(hash.data, hash.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_Password(void* ptr)
{
	return ({ QByteArray* t833eaa = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->password().toUtf8()); Moc_PackedString { const_cast<char*>(t833eaa->prepend("WHITESPACE").constData()+10), t833eaa->size()-10, t833eaa }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_PasswordDefault(void* ptr)
{
	return ({ QByteArray* tbe2949 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->passwordDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tbe2949->prepend("WHITESPACE").constData()+10), tbe2949->size()-10, tbe2949 }; });
}

void ApproveSignDataCtx687eda_SetPassword(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setPassword(QString::fromUtf8(password.data, password.len));
}

void ApproveSignDataCtx687eda_SetPasswordDefault(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setPasswordDefault(QString::fromUtf8(password.data, password.len));
}

void ApproveSignDataCtx687eda_ConnectPasswordChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::passwordChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_PasswordChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectPasswordChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::passwordChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_PasswordChanged));
}

void ApproveSignDataCtx687eda_PasswordChanged(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->passwordChanged(QString::fromUtf8(password.data, password.len));
}

struct Moc_PackedString ApproveSignDataCtx687eda_FromSrc(void* ptr)
{
	return ({ QByteArray* tc808ad = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->fromSrc().toUtf8()); Moc_PackedString { const_cast<char*>(tc808ad->prepend("WHITESPACE").constData()+10), tc808ad->size()-10, tc808ad }; });
}

struct Moc_PackedString ApproveSignDataCtx687eda_FromSrcDefault(void* ptr)
{
	return ({ QByteArray* t981347 = new QByteArray(static_cast<ApproveSignDataCtx687eda*>(ptr)->fromSrcDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t981347->prepend("WHITESPACE").constData()+10), t981347->size()-10, t981347 }; });
}

void ApproveSignDataCtx687eda_SetFromSrc(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setFromSrc(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

void ApproveSignDataCtx687eda_SetFromSrcDefault(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->setFromSrcDefault(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

void ApproveSignDataCtx687eda_ConnectFromSrcChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::fromSrcChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_FromSrcChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveSignDataCtx687eda_DisconnectFromSrcChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::fromSrcChanged), static_cast<ApproveSignDataCtx687eda*>(ptr), static_cast<void (ApproveSignDataCtx687eda::*)(QString)>(&ApproveSignDataCtx687eda::Signal_FromSrcChanged));
}

void ApproveSignDataCtx687eda_FromSrcChanged(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->fromSrcChanged(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

int ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType()
{
	return qRegisterMetaType<ApproveSignDataCtx687eda*>();
}

int ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<ApproveSignDataCtx687eda*>(const_cast<const char*>(typeName));
}

int ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveSignDataCtx687eda>();
#else
	return 0;
#endif
}

int ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveSignDataCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<ApproveSignDataCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* ApproveSignDataCtx687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveSignDataCtx687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveSignDataCtx687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* ApproveSignDataCtx687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void ApproveSignDataCtx687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* ApproveSignDataCtx687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* ApproveSignDataCtx687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveSignDataCtx687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveSignDataCtx687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveSignDataCtx687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveSignDataCtx687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveSignDataCtx687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveSignDataCtx687eda_NewApproveSignDataCtx(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveSignDataCtx687eda(static_cast<QWindow*>(parent));
	} else {
		return new ApproveSignDataCtx687eda(static_cast<QObject*>(parent));
	}
}

void ApproveSignDataCtx687eda_DestroyApproveSignDataCtx(void* ptr)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->~ApproveSignDataCtx687eda();
}

void ApproveSignDataCtx687eda_DestroyApproveSignDataCtxDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void ApproveSignDataCtx687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void ApproveSignDataCtx687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void ApproveSignDataCtx687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void ApproveSignDataCtx687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::deleteLater();
}

void ApproveSignDataCtx687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char ApproveSignDataCtx687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char ApproveSignDataCtx687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void ApproveSignDataCtx687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<ApproveSignDataCtx687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void ApproveTxCtx687eda_ConnectApprove(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::approve), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Approve), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectApprove(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::approve), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Approve));
}

void ApproveTxCtx687eda_Approve(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->approve();
}

void ApproveTxCtx687eda_ConnectReject(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::reject), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Reject), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectReject(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::reject), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Reject));
}

void ApproveTxCtx687eda_Reject(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->reject();
}

void ApproveTxCtx687eda_ConnectCheckTxDiff(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::checkTxDiff), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_CheckTxDiff), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectCheckTxDiff(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::checkTxDiff), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_CheckTxDiff));
}

void ApproveTxCtx687eda_CheckTxDiff(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->checkTxDiff();
}

void ApproveTxCtx687eda_ConnectBack(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::back), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Back), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectBack(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::back), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)()>(&ApproveTxCtx687eda::Signal_Back));
}

void ApproveTxCtx687eda_Back(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->back();
}

void ApproveTxCtx687eda_ConnectEdited(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString, QString)>(&ApproveTxCtx687eda::edited), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString, QString)>(&ApproveTxCtx687eda::Signal_Edited), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectEdited(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString, QString)>(&ApproveTxCtx687eda::edited), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString, QString)>(&ApproveTxCtx687eda::Signal_Edited));
}

void ApproveTxCtx687eda_Edited(void* ptr, struct Moc_PackedString s, struct Moc_PackedString v)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->edited(QString::fromUtf8(s.data, s.len), QString::fromUtf8(v.data, v.len));
}

void ApproveTxCtx687eda_ConnectChangeValueUnit(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::changeValueUnit), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ChangeValueUnit), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectChangeValueUnit(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::changeValueUnit), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ChangeValueUnit));
}

void ApproveTxCtx687eda_ChangeValueUnit(void* ptr, int v)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->changeValueUnit(v);
}

void ApproveTxCtx687eda_ConnectChangeGasPriceUnit(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::changeGasPriceUnit), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ChangeGasPriceUnit), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectChangeGasPriceUnit(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::changeGasPriceUnit), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ChangeGasPriceUnit));
}

void ApproveTxCtx687eda_ChangeGasPriceUnit(void* ptr, int v)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->changeGasPriceUnit(v);
}

int ApproveTxCtx687eda_ValueUnit(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->valueUnit();
}

int ApproveTxCtx687eda_ValueUnitDefault(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->valueUnitDefault();
}

void ApproveTxCtx687eda_SetValueUnit(void* ptr, int valueUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setValueUnit(valueUnit);
}

void ApproveTxCtx687eda_SetValueUnitDefault(void* ptr, int valueUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setValueUnitDefault(valueUnit);
}

void ApproveTxCtx687eda_ConnectValueUnitChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::valueUnitChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ValueUnitChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectValueUnitChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::valueUnitChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_ValueUnitChanged));
}

void ApproveTxCtx687eda_ValueUnitChanged(void* ptr, int valueUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->valueUnitChanged(valueUnit);
}

struct Moc_PackedString ApproveTxCtx687eda_Remote(void* ptr)
{
	return ({ QByteArray* t90f7f1 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->remote().toUtf8()); Moc_PackedString { const_cast<char*>(t90f7f1->prepend("WHITESPACE").constData()+10), t90f7f1->size()-10, t90f7f1 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_RemoteDefault(void* ptr)
{
	return ({ QByteArray* t5fca0a = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->remoteDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t5fca0a->prepend("WHITESPACE").constData()+10), t5fca0a->size()-10, t5fca0a }; });
}

void ApproveTxCtx687eda_SetRemote(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setRemote(QString::fromUtf8(remote.data, remote.len));
}

void ApproveTxCtx687eda_SetRemoteDefault(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setRemoteDefault(QString::fromUtf8(remote.data, remote.len));
}

void ApproveTxCtx687eda_ConnectRemoteChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::remoteChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_RemoteChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectRemoteChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::remoteChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_RemoteChanged));
}

void ApproveTxCtx687eda_RemoteChanged(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->remoteChanged(QString::fromUtf8(remote.data, remote.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Transport(void* ptr)
{
	return ({ QByteArray* te6d44f = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->transport().toUtf8()); Moc_PackedString { const_cast<char*>(te6d44f->prepend("WHITESPACE").constData()+10), te6d44f->size()-10, te6d44f }; });
}

struct Moc_PackedString ApproveTxCtx687eda_TransportDefault(void* ptr)
{
	return ({ QByteArray* t09c1b3 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->transportDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t09c1b3->prepend("WHITESPACE").constData()+10), t09c1b3->size()-10, t09c1b3 }; });
}

void ApproveTxCtx687eda_SetTransport(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setTransport(QString::fromUtf8(transport.data, transport.len));
}

void ApproveTxCtx687eda_SetTransportDefault(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setTransportDefault(QString::fromUtf8(transport.data, transport.len));
}

void ApproveTxCtx687eda_ConnectTransportChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::transportChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_TransportChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectTransportChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::transportChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_TransportChanged));
}

void ApproveTxCtx687eda_TransportChanged(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->transportChanged(QString::fromUtf8(transport.data, transport.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Endpoint(void* ptr)
{
	return ({ QByteArray* taa8e81 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->endpoint().toUtf8()); Moc_PackedString { const_cast<char*>(taa8e81->prepend("WHITESPACE").constData()+10), taa8e81->size()-10, taa8e81 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_EndpointDefault(void* ptr)
{
	return ({ QByteArray* t2b1328 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->endpointDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t2b1328->prepend("WHITESPACE").constData()+10), t2b1328->size()-10, t2b1328 }; });
}

void ApproveTxCtx687eda_SetEndpoint(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setEndpoint(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveTxCtx687eda_SetEndpointDefault(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setEndpointDefault(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveTxCtx687eda_ConnectEndpointChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::endpointChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_EndpointChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectEndpointChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::endpointChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_EndpointChanged));
}

void ApproveTxCtx687eda_EndpointChanged(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->endpointChanged(QString::fromUtf8(endpoint.data, endpoint.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Data(void* ptr)
{
	return ({ QByteArray* t8072b7 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->data().toUtf8()); Moc_PackedString { const_cast<char*>(t8072b7->prepend("WHITESPACE").constData()+10), t8072b7->size()-10, t8072b7 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_DataDefault(void* ptr)
{
	return ({ QByteArray* t04fe35 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->dataDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t04fe35->prepend("WHITESPACE").constData()+10), t04fe35->size()-10, t04fe35 }; });
}

void ApproveTxCtx687eda_SetData(void* ptr, struct Moc_PackedString data)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setData(QString::fromUtf8(data.data, data.len));
}

void ApproveTxCtx687eda_SetDataDefault(void* ptr, struct Moc_PackedString data)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setDataDefault(QString::fromUtf8(data.data, data.len));
}

void ApproveTxCtx687eda_ConnectDataChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::dataChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_DataChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectDataChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::dataChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_DataChanged));
}

void ApproveTxCtx687eda_DataChanged(void* ptr, struct Moc_PackedString data)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->dataChanged(QString::fromUtf8(data.data, data.len));
}

struct Moc_PackedString ApproveTxCtx687eda_From(void* ptr)
{
	return ({ QByteArray* tc27f9b = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->from().toUtf8()); Moc_PackedString { const_cast<char*>(tc27f9b->prepend("WHITESPACE").constData()+10), tc27f9b->size()-10, tc27f9b }; });
}

struct Moc_PackedString ApproveTxCtx687eda_FromDefault(void* ptr)
{
	return ({ QByteArray* t805f9b = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->fromDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t805f9b->prepend("WHITESPACE").constData()+10), t805f9b->size()-10, t805f9b }; });
}

void ApproveTxCtx687eda_SetFrom(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFrom(QString::fromUtf8(from.data, from.len));
}

void ApproveTxCtx687eda_SetFromDefault(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromDefault(QString::fromUtf8(from.data, from.len));
}

void ApproveTxCtx687eda_ConnectFromChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectFromChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromChanged));
}

void ApproveTxCtx687eda_FromChanged(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->fromChanged(QString::fromUtf8(from.data, from.len));
}

struct Moc_PackedString ApproveTxCtx687eda_FromWarning(void* ptr)
{
	return ({ QByteArray* td56e55 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->fromWarning().toUtf8()); Moc_PackedString { const_cast<char*>(td56e55->prepend("WHITESPACE").constData()+10), td56e55->size()-10, td56e55 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_FromWarningDefault(void* ptr)
{
	return ({ QByteArray* t2020a1 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->fromWarningDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t2020a1->prepend("WHITESPACE").constData()+10), t2020a1->size()-10, t2020a1 }; });
}

void ApproveTxCtx687eda_SetFromWarning(void* ptr, struct Moc_PackedString fromWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromWarning(QString::fromUtf8(fromWarning.data, fromWarning.len));
}

void ApproveTxCtx687eda_SetFromWarningDefault(void* ptr, struct Moc_PackedString fromWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromWarningDefault(QString::fromUtf8(fromWarning.data, fromWarning.len));
}

void ApproveTxCtx687eda_ConnectFromWarningChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromWarningChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromWarningChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectFromWarningChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromWarningChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromWarningChanged));
}

void ApproveTxCtx687eda_FromWarningChanged(void* ptr, struct Moc_PackedString fromWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->fromWarningChanged(QString::fromUtf8(fromWarning.data, fromWarning.len));
}

char ApproveTxCtx687eda_IsFromVisible(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->isFromVisible();
}

char ApproveTxCtx687eda_IsFromVisibleDefault(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->isFromVisibleDefault();
}

void ApproveTxCtx687eda_SetFromVisible(void* ptr, char fromVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromVisible(fromVisible != 0);
}

void ApproveTxCtx687eda_SetFromVisibleDefault(void* ptr, char fromVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromVisibleDefault(fromVisible != 0);
}

void ApproveTxCtx687eda_ConnectFromVisibleChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::fromVisibleChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::Signal_FromVisibleChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectFromVisibleChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::fromVisibleChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::Signal_FromVisibleChanged));
}

void ApproveTxCtx687eda_FromVisibleChanged(void* ptr, char fromVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->fromVisibleChanged(fromVisible != 0);
}

struct Moc_PackedString ApproveTxCtx687eda_To(void* ptr)
{
	return ({ QByteArray* t6e4b06 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->to().toUtf8()); Moc_PackedString { const_cast<char*>(t6e4b06->prepend("WHITESPACE").constData()+10), t6e4b06->size()-10, t6e4b06 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_ToDefault(void* ptr)
{
	return ({ QByteArray* t7ceec8 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->toDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t7ceec8->prepend("WHITESPACE").constData()+10), t7ceec8->size()-10, t7ceec8 }; });
}

void ApproveTxCtx687eda_SetTo(void* ptr, struct Moc_PackedString to)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setTo(QString::fromUtf8(to.data, to.len));
}

void ApproveTxCtx687eda_SetToDefault(void* ptr, struct Moc_PackedString to)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToDefault(QString::fromUtf8(to.data, to.len));
}

void ApproveTxCtx687eda_ConnectToChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectToChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToChanged));
}

void ApproveTxCtx687eda_ToChanged(void* ptr, struct Moc_PackedString to)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->toChanged(QString::fromUtf8(to.data, to.len));
}

struct Moc_PackedString ApproveTxCtx687eda_ToWarning(void* ptr)
{
	return ({ QByteArray* t7006fd = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->toWarning().toUtf8()); Moc_PackedString { const_cast<char*>(t7006fd->prepend("WHITESPACE").constData()+10), t7006fd->size()-10, t7006fd }; });
}

struct Moc_PackedString ApproveTxCtx687eda_ToWarningDefault(void* ptr)
{
	return ({ QByteArray* t523608 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->toWarningDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t523608->prepend("WHITESPACE").constData()+10), t523608->size()-10, t523608 }; });
}

void ApproveTxCtx687eda_SetToWarning(void* ptr, struct Moc_PackedString toWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToWarning(QString::fromUtf8(toWarning.data, toWarning.len));
}

void ApproveTxCtx687eda_SetToWarningDefault(void* ptr, struct Moc_PackedString toWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToWarningDefault(QString::fromUtf8(toWarning.data, toWarning.len));
}

void ApproveTxCtx687eda_ConnectToWarningChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toWarningChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToWarningChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectToWarningChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toWarningChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToWarningChanged));
}

void ApproveTxCtx687eda_ToWarningChanged(void* ptr, struct Moc_PackedString toWarning)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->toWarningChanged(QString::fromUtf8(toWarning.data, toWarning.len));
}

char ApproveTxCtx687eda_IsToVisible(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->isToVisible();
}

char ApproveTxCtx687eda_IsToVisibleDefault(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->isToVisibleDefault();
}

void ApproveTxCtx687eda_SetToVisible(void* ptr, char toVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToVisible(toVisible != 0);
}

void ApproveTxCtx687eda_SetToVisibleDefault(void* ptr, char toVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToVisibleDefault(toVisible != 0);
}

void ApproveTxCtx687eda_ConnectToVisibleChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::toVisibleChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::Signal_ToVisibleChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectToVisibleChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::toVisibleChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(bool)>(&ApproveTxCtx687eda::Signal_ToVisibleChanged));
}

void ApproveTxCtx687eda_ToVisibleChanged(void* ptr, char toVisible)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->toVisibleChanged(toVisible != 0);
}

struct Moc_PackedString ApproveTxCtx687eda_Gas(void* ptr)
{
	return ({ QByteArray* t3b252a = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->gas().toUtf8()); Moc_PackedString { const_cast<char*>(t3b252a->prepend("WHITESPACE").constData()+10), t3b252a->size()-10, t3b252a }; });
}

struct Moc_PackedString ApproveTxCtx687eda_GasDefault(void* ptr)
{
	return ({ QByteArray* tdee409 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->gasDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tdee409->prepend("WHITESPACE").constData()+10), tdee409->size()-10, tdee409 }; });
}

void ApproveTxCtx687eda_SetGas(void* ptr, struct Moc_PackedString gas)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGas(QString::fromUtf8(gas.data, gas.len));
}

void ApproveTxCtx687eda_SetGasDefault(void* ptr, struct Moc_PackedString gas)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGasDefault(QString::fromUtf8(gas.data, gas.len));
}

void ApproveTxCtx687eda_ConnectGasChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::gasChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_GasChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectGasChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::gasChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_GasChanged));
}

void ApproveTxCtx687eda_GasChanged(void* ptr, struct Moc_PackedString gas)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->gasChanged(QString::fromUtf8(gas.data, gas.len));
}

struct Moc_PackedString ApproveTxCtx687eda_GasPrice(void* ptr)
{
	return ({ QByteArray* te19d41 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->gasPrice().toUtf8()); Moc_PackedString { const_cast<char*>(te19d41->prepend("WHITESPACE").constData()+10), te19d41->size()-10, te19d41 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_GasPriceDefault(void* ptr)
{
	return ({ QByteArray* tef8e7f = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->gasPriceDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tef8e7f->prepend("WHITESPACE").constData()+10), tef8e7f->size()-10, tef8e7f }; });
}

void ApproveTxCtx687eda_SetGasPrice(void* ptr, struct Moc_PackedString gasPrice)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGasPrice(QString::fromUtf8(gasPrice.data, gasPrice.len));
}

void ApproveTxCtx687eda_SetGasPriceDefault(void* ptr, struct Moc_PackedString gasPrice)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGasPriceDefault(QString::fromUtf8(gasPrice.data, gasPrice.len));
}

void ApproveTxCtx687eda_ConnectGasPriceChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::gasPriceChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_GasPriceChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectGasPriceChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::gasPriceChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_GasPriceChanged));
}

void ApproveTxCtx687eda_GasPriceChanged(void* ptr, struct Moc_PackedString gasPrice)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->gasPriceChanged(QString::fromUtf8(gasPrice.data, gasPrice.len));
}

int ApproveTxCtx687eda_GasPriceUnit(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->gasPriceUnit();
}

int ApproveTxCtx687eda_GasPriceUnitDefault(void* ptr)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->gasPriceUnitDefault();
}

void ApproveTxCtx687eda_SetGasPriceUnit(void* ptr, int gasPriceUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGasPriceUnit(gasPriceUnit);
}

void ApproveTxCtx687eda_SetGasPriceUnitDefault(void* ptr, int gasPriceUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setGasPriceUnitDefault(gasPriceUnit);
}

void ApproveTxCtx687eda_ConnectGasPriceUnitChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::gasPriceUnitChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_GasPriceUnitChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectGasPriceUnitChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::gasPriceUnitChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(qint32)>(&ApproveTxCtx687eda::Signal_GasPriceUnitChanged));
}

void ApproveTxCtx687eda_GasPriceUnitChanged(void* ptr, int gasPriceUnit)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->gasPriceUnitChanged(gasPriceUnit);
}

struct Moc_PackedString ApproveTxCtx687eda_Nonce(void* ptr)
{
	return ({ QByteArray* t532f58 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->nonce().toUtf8()); Moc_PackedString { const_cast<char*>(t532f58->prepend("WHITESPACE").constData()+10), t532f58->size()-10, t532f58 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_NonceDefault(void* ptr)
{
	return ({ QByteArray* tf1f93d = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->nonceDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tf1f93d->prepend("WHITESPACE").constData()+10), tf1f93d->size()-10, tf1f93d }; });
}

void ApproveTxCtx687eda_SetNonce(void* ptr, struct Moc_PackedString nonce)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setNonce(QString::fromUtf8(nonce.data, nonce.len));
}

void ApproveTxCtx687eda_SetNonceDefault(void* ptr, struct Moc_PackedString nonce)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setNonceDefault(QString::fromUtf8(nonce.data, nonce.len));
}

void ApproveTxCtx687eda_ConnectNonceChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::nonceChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_NonceChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectNonceChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::nonceChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_NonceChanged));
}

void ApproveTxCtx687eda_NonceChanged(void* ptr, struct Moc_PackedString nonce)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->nonceChanged(QString::fromUtf8(nonce.data, nonce.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Value(void* ptr)
{
	return ({ QByteArray* te5a3fa = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->value().toUtf8()); Moc_PackedString { const_cast<char*>(te5a3fa->prepend("WHITESPACE").constData()+10), te5a3fa->size()-10, te5a3fa }; });
}

struct Moc_PackedString ApproveTxCtx687eda_ValueDefault(void* ptr)
{
	return ({ QByteArray* t9772be = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->valueDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t9772be->prepend("WHITESPACE").constData()+10), t9772be->size()-10, t9772be }; });
}

void ApproveTxCtx687eda_SetValue(void* ptr, struct Moc_PackedString value)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setValue(QString::fromUtf8(value.data, value.len));
}

void ApproveTxCtx687eda_SetValueDefault(void* ptr, struct Moc_PackedString value)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setValueDefault(QString::fromUtf8(value.data, value.len));
}

void ApproveTxCtx687eda_ConnectValueChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::valueChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ValueChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectValueChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::valueChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ValueChanged));
}

void ApproveTxCtx687eda_ValueChanged(void* ptr, struct Moc_PackedString value)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->valueChanged(QString::fromUtf8(value.data, value.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Password(void* ptr)
{
	return ({ QByteArray* t918077 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->password().toUtf8()); Moc_PackedString { const_cast<char*>(t918077->prepend("WHITESPACE").constData()+10), t918077->size()-10, t918077 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_PasswordDefault(void* ptr)
{
	return ({ QByteArray* t686f68 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->passwordDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t686f68->prepend("WHITESPACE").constData()+10), t686f68->size()-10, t686f68 }; });
}

void ApproveTxCtx687eda_SetPassword(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setPassword(QString::fromUtf8(password.data, password.len));
}

void ApproveTxCtx687eda_SetPasswordDefault(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setPasswordDefault(QString::fromUtf8(password.data, password.len));
}

void ApproveTxCtx687eda_ConnectPasswordChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::passwordChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_PasswordChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectPasswordChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::passwordChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_PasswordChanged));
}

void ApproveTxCtx687eda_PasswordChanged(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->passwordChanged(QString::fromUtf8(password.data, password.len));
}

struct Moc_PackedString ApproveTxCtx687eda_FromSrc(void* ptr)
{
	return ({ QByteArray* t6aa1b7 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->fromSrc().toUtf8()); Moc_PackedString { const_cast<char*>(t6aa1b7->prepend("WHITESPACE").constData()+10), t6aa1b7->size()-10, t6aa1b7 }; });
}

struct Moc_PackedString ApproveTxCtx687eda_FromSrcDefault(void* ptr)
{
	return ({ QByteArray* tc7a000 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->fromSrcDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tc7a000->prepend("WHITESPACE").constData()+10), tc7a000->size()-10, tc7a000 }; });
}

void ApproveTxCtx687eda_SetFromSrc(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromSrc(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

void ApproveTxCtx687eda_SetFromSrcDefault(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setFromSrcDefault(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

void ApproveTxCtx687eda_ConnectFromSrcChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromSrcChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromSrcChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectFromSrcChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::fromSrcChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_FromSrcChanged));
}

void ApproveTxCtx687eda_FromSrcChanged(void* ptr, struct Moc_PackedString fromSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->fromSrcChanged(QString::fromUtf8(fromSrc.data, fromSrc.len));
}

struct Moc_PackedString ApproveTxCtx687eda_ToSrc(void* ptr)
{
	return ({ QByteArray* t03230d = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->toSrc().toUtf8()); Moc_PackedString { const_cast<char*>(t03230d->prepend("WHITESPACE").constData()+10), t03230d->size()-10, t03230d }; });
}

struct Moc_PackedString ApproveTxCtx687eda_ToSrcDefault(void* ptr)
{
	return ({ QByteArray* t20ab85 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->toSrcDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t20ab85->prepend("WHITESPACE").constData()+10), t20ab85->size()-10, t20ab85 }; });
}

void ApproveTxCtx687eda_SetToSrc(void* ptr, struct Moc_PackedString toSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToSrc(QString::fromUtf8(toSrc.data, toSrc.len));
}

void ApproveTxCtx687eda_SetToSrcDefault(void* ptr, struct Moc_PackedString toSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setToSrcDefault(QString::fromUtf8(toSrc.data, toSrc.len));
}

void ApproveTxCtx687eda_ConnectToSrcChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toSrcChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToSrcChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectToSrcChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::toSrcChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_ToSrcChanged));
}

void ApproveTxCtx687eda_ToSrcChanged(void* ptr, struct Moc_PackedString toSrc)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->toSrcChanged(QString::fromUtf8(toSrc.data, toSrc.len));
}

struct Moc_PackedString ApproveTxCtx687eda_Diff(void* ptr)
{
	return ({ QByteArray* t81763e = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->diff().toUtf8()); Moc_PackedString { const_cast<char*>(t81763e->prepend("WHITESPACE").constData()+10), t81763e->size()-10, t81763e }; });
}

struct Moc_PackedString ApproveTxCtx687eda_DiffDefault(void* ptr)
{
	return ({ QByteArray* td180f6 = new QByteArray(static_cast<ApproveTxCtx687eda*>(ptr)->diffDefault().toUtf8()); Moc_PackedString { const_cast<char*>(td180f6->prepend("WHITESPACE").constData()+10), td180f6->size()-10, td180f6 }; });
}

void ApproveTxCtx687eda_SetDiff(void* ptr, struct Moc_PackedString diff)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setDiff(QString::fromUtf8(diff.data, diff.len));
}

void ApproveTxCtx687eda_SetDiffDefault(void* ptr, struct Moc_PackedString diff)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->setDiffDefault(QString::fromUtf8(diff.data, diff.len));
}

void ApproveTxCtx687eda_ConnectDiffChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::diffChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_DiffChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveTxCtx687eda_DisconnectDiffChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::diffChanged), static_cast<ApproveTxCtx687eda*>(ptr), static_cast<void (ApproveTxCtx687eda::*)(QString)>(&ApproveTxCtx687eda::Signal_DiffChanged));
}

void ApproveTxCtx687eda_DiffChanged(void* ptr, struct Moc_PackedString diff)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->diffChanged(QString::fromUtf8(diff.data, diff.len));
}

int ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType()
{
	return qRegisterMetaType<ApproveTxCtx687eda*>();
}

int ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<ApproveTxCtx687eda*>(const_cast<const char*>(typeName));
}

int ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveTxCtx687eda>();
#else
	return 0;
#endif
}

int ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveTxCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<ApproveTxCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* ApproveTxCtx687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveTxCtx687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveTxCtx687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* ApproveTxCtx687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void ApproveTxCtx687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* ApproveTxCtx687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* ApproveTxCtx687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveTxCtx687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveTxCtx687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveTxCtx687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveTxCtx687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveTxCtx687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveTxCtx687eda_NewApproveTxCtx(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveTxCtx687eda(static_cast<QWindow*>(parent));
	} else {
		return new ApproveTxCtx687eda(static_cast<QObject*>(parent));
	}
}

void ApproveTxCtx687eda_DestroyApproveTxCtx(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->~ApproveTxCtx687eda();
}

void ApproveTxCtx687eda_DestroyApproveTxCtxDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void ApproveTxCtx687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void ApproveTxCtx687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void ApproveTxCtx687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void ApproveTxCtx687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::deleteLater();
}

void ApproveTxCtx687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char ApproveTxCtx687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char ApproveTxCtx687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<ApproveTxCtx687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void ApproveTxCtx687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<ApproveTxCtx687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void TxListAccountsModel687eda_ConnectAdd(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListAccountsModel687eda*>(ptr), static_cast<void (TxListAccountsModel687eda::*)(QString)>(&TxListAccountsModel687eda::add), static_cast<TxListAccountsModel687eda*>(ptr), static_cast<void (TxListAccountsModel687eda::*)(QString)>(&TxListAccountsModel687eda::Signal_Add), static_cast<Qt::ConnectionType>(t));
}

void TxListAccountsModel687eda_DisconnectAdd(void* ptr)
{
	QObject::disconnect(static_cast<TxListAccountsModel687eda*>(ptr), static_cast<void (TxListAccountsModel687eda::*)(QString)>(&TxListAccountsModel687eda::add), static_cast<TxListAccountsModel687eda*>(ptr), static_cast<void (TxListAccountsModel687eda::*)(QString)>(&TxListAccountsModel687eda::Signal_Add));
}

void TxListAccountsModel687eda_Add(void* ptr, struct Moc_PackedString tx)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->add(QString::fromUtf8(tx.data, tx.len));
}

int TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType()
{
	return qRegisterMetaType<TxListAccountsModel687eda*>();
}

int TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<TxListAccountsModel687eda*>(const_cast<const char*>(typeName));
}

int TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListAccountsModel687eda>();
#else
	return 0;
#endif
}

int TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListAccountsModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<TxListAccountsModel687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

int TxListAccountsModel687eda_____itemData_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda_____itemData_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda_____itemData_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListAccountsModel687eda_____roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda_____roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda_____roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListAccountsModel687eda_____setItemData_roles_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda_____setItemData_roles_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda_____setItemData_roles_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* TxListAccountsModel687eda___changePersistentIndexList_from_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___changePersistentIndexList_from_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListAccountsModel687eda___changePersistentIndexList_from_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListAccountsModel687eda___changePersistentIndexList_to_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___changePersistentIndexList_to_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListAccountsModel687eda___changePersistentIndexList_to_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

int TxListAccountsModel687eda___dataChanged_roles_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QVector<int>*>(ptr)->at(i); if (i == static_cast<QVector<int>*>(ptr)->size()-1) { static_cast<QVector<int>*>(ptr)->~QVector(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda___dataChanged_roles_setList(void* ptr, int i)
{
	static_cast<QVector<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda___dataChanged_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QVector<int>();
}

void* TxListAccountsModel687eda___itemData_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___itemData_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* TxListAccountsModel687eda___itemData_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList TxListAccountsModel687eda___itemData_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

void* TxListAccountsModel687eda___layoutAboutToBeChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___layoutAboutToBeChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* TxListAccountsModel687eda___layoutAboutToBeChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* TxListAccountsModel687eda___layoutChanged_parents_atList(void* ptr, int i)
{
	return new QPersistentModelIndex(({QPersistentModelIndex tmp = static_cast<QList<QPersistentModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QPersistentModelIndex>*>(ptr)->size()-1) { static_cast<QList<QPersistentModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___layoutChanged_parents_setList(void* ptr, void* i)
{
	static_cast<QList<QPersistentModelIndex>*>(ptr)->append(*static_cast<QPersistentModelIndex*>(i));
}

void* TxListAccountsModel687eda___layoutChanged_parents_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QPersistentModelIndex>();
}

void* TxListAccountsModel687eda___match_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___match_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListAccountsModel687eda___match_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListAccountsModel687eda___mimeData_indexes_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___mimeData_indexes_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListAccountsModel687eda___mimeData_indexes_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListAccountsModel687eda___persistentIndexList_atList(void* ptr, int i)
{
	return new QModelIndex(({QModelIndex tmp = static_cast<QList<QModelIndex>*>(ptr)->at(i); if (i == static_cast<QList<QModelIndex>*>(ptr)->size()-1) { static_cast<QList<QModelIndex>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___persistentIndexList_setList(void* ptr, void* i)
{
	static_cast<QList<QModelIndex>*>(ptr)->append(*static_cast<QModelIndex*>(i));
}

void* TxListAccountsModel687eda___persistentIndexList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QModelIndex>();
}

void* TxListAccountsModel687eda___roleNames_atList(void* ptr, int v, int i)
{
	return new QByteArray(({ QByteArray tmp = static_cast<QHash<int, QByteArray>*>(ptr)->value(v); if (i == static_cast<QHash<int, QByteArray>*>(ptr)->size()-1) { static_cast<QHash<int, QByteArray>*>(ptr)->~QHash(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___roleNames_setList(void* ptr, int key, void* i)
{
	static_cast<QHash<int, QByteArray>*>(ptr)->insert(key, *static_cast<QByteArray*>(i));
}

void* TxListAccountsModel687eda___roleNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QHash<int, QByteArray>();
}

struct Moc_PackedList TxListAccountsModel687eda___roleNames_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue7fc3bb = new QList<int>(static_cast<QHash<int, QByteArray>*>(ptr)->keys()); Moc_PackedList { tmpValue7fc3bb, tmpValue7fc3bb->size() }; });
}

void* TxListAccountsModel687eda___setItemData_roles_atList(void* ptr, int v, int i)
{
	return new QVariant(({ QVariant tmp = static_cast<QMap<int, QVariant>*>(ptr)->value(v); if (i == static_cast<QMap<int, QVariant>*>(ptr)->size()-1) { static_cast<QMap<int, QVariant>*>(ptr)->~QMap(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___setItemData_roles_setList(void* ptr, int key, void* i)
{
	static_cast<QMap<int, QVariant>*>(ptr)->insert(key, *static_cast<QVariant*>(i));
}

void* TxListAccountsModel687eda___setItemData_roles_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QMap<int, QVariant>();
}

struct Moc_PackedList TxListAccountsModel687eda___setItemData_roles_keyList(void* ptr)
{
	return ({ QList<int>* tmpValue249128 = new QList<int>(static_cast<QMap<int, QVariant>*>(ptr)->keys()); Moc_PackedList { tmpValue249128, tmpValue249128->size() }; });
}

int TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

int TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_atList(void* ptr, int i)
{
	return ({int tmp = static_cast<QList<int>*>(ptr)->at(i); if (i == static_cast<QList<int>*>(ptr)->size()-1) { static_cast<QList<int>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_setList(void* ptr, int i)
{
	static_cast<QList<int>*>(ptr)->append(i);
}

void* TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<int>();
}

void* TxListAccountsModel687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListAccountsModel687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* TxListAccountsModel687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListAccountsModel687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* TxListAccountsModel687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* TxListAccountsModel687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListAccountsModel687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListAccountsModel687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListAccountsModel687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListAccountsModel687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListAccountsModel687eda_NewTxListAccountsModel(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new TxListAccountsModel687eda(static_cast<QWindow*>(parent));
	} else {
		return new TxListAccountsModel687eda(static_cast<QObject*>(parent));
	}
}

void TxListAccountsModel687eda_DestroyTxListAccountsModel(void* ptr)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->~TxListAccountsModel687eda();
}

void TxListAccountsModel687eda_DestroyTxListAccountsModelDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

char TxListAccountsModel687eda_DropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::dropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

long long TxListAccountsModel687eda_FlagsDefault(void* ptr, void* index)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::flags(*static_cast<QModelIndex*>(index));
}

void* TxListAccountsModel687eda_IndexDefault(void* ptr, int row, int column, void* parent)
{
	return new QModelIndex(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::index(row, column, *static_cast<QModelIndex*>(parent)));
}

void* TxListAccountsModel687eda_SiblingDefault(void* ptr, int row, int column, void* idx)
{
	return new QModelIndex(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::sibling(row, column, *static_cast<QModelIndex*>(idx)));
}

void* TxListAccountsModel687eda_BuddyDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::buddy(*static_cast<QModelIndex*>(index)));
}

char TxListAccountsModel687eda_CanDropMimeDataDefault(void* ptr, void* data, long long action, int row, int column, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::canDropMimeData(static_cast<QMimeData*>(data), static_cast<Qt::DropAction>(action), row, column, *static_cast<QModelIndex*>(parent));
}

char TxListAccountsModel687eda_CanFetchMoreDefault(void* ptr, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::canFetchMore(*static_cast<QModelIndex*>(parent));
}

int TxListAccountsModel687eda_ColumnCountDefault(void* ptr, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::columnCount(*static_cast<QModelIndex*>(parent));
}

void* TxListAccountsModel687eda_DataDefault(void* ptr, void* index, int role)
{
	Q_UNUSED(ptr);
	Q_UNUSED(index);
	Q_UNUSED(role);

}

void TxListAccountsModel687eda_FetchMoreDefault(void* ptr, void* parent)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::fetchMore(*static_cast<QModelIndex*>(parent));
}

char TxListAccountsModel687eda_HasChildrenDefault(void* ptr, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::hasChildren(*static_cast<QModelIndex*>(parent));
}

void* TxListAccountsModel687eda_HeaderDataDefault(void* ptr, int section, long long orientation, int role)
{
	return new QVariant(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::headerData(section, static_cast<Qt::Orientation>(orientation), role));
}

char TxListAccountsModel687eda_InsertColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::insertColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char TxListAccountsModel687eda_InsertRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::insertRows(row, count, *static_cast<QModelIndex*>(parent));
}

struct Moc_PackedList TxListAccountsModel687eda_ItemDataDefault(void* ptr, void* index)
{
	return ({ QMap<int, QVariant>* tmpValue4772a3 = new QMap<int, QVariant>(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::itemData(*static_cast<QModelIndex*>(index))); Moc_PackedList { tmpValue4772a3, tmpValue4772a3->size() }; });
}

struct Moc_PackedList TxListAccountsModel687eda_MatchDefault(void* ptr, void* start, int role, void* value, int hits, long long flags)
{
	return ({ QList<QModelIndex>* tmpValued34c09 = new QList<QModelIndex>(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::match(*static_cast<QModelIndex*>(start), role, *static_cast<QVariant*>(value), hits, static_cast<Qt::MatchFlag>(flags))); Moc_PackedList { tmpValued34c09, tmpValued34c09->size() }; });
}

void* TxListAccountsModel687eda_MimeDataDefault(void* ptr, void* indexes)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::mimeData(({ QList<QModelIndex>* tmpP = static_cast<QList<QModelIndex>*>(indexes); QList<QModelIndex> tmpV = *tmpP; tmpP->~QList(); free(tmpP); tmpV; }));
}

struct Moc_PackedString TxListAccountsModel687eda_MimeTypesDefault(void* ptr)
{
	return ({ QByteArray* tb2f4e1 = new QByteArray(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::mimeTypes().join("¡¦!").toUtf8()); Moc_PackedString { const_cast<char*>(tb2f4e1->prepend("WHITESPACE").constData()+10), tb2f4e1->size()-10, tb2f4e1 }; });
}

char TxListAccountsModel687eda_MoveColumnsDefault(void* ptr, void* sourceParent, int sourceColumn, int count, void* destinationParent, int destinationChild)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::moveColumns(*static_cast<QModelIndex*>(sourceParent), sourceColumn, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

char TxListAccountsModel687eda_MoveRowsDefault(void* ptr, void* sourceParent, int sourceRow, int count, void* destinationParent, int destinationChild)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::moveRows(*static_cast<QModelIndex*>(sourceParent), sourceRow, count, *static_cast<QModelIndex*>(destinationParent), destinationChild);
}

void* TxListAccountsModel687eda_ParentDefault(void* ptr, void* index)
{
	return new QModelIndex(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::parent(*static_cast<QModelIndex*>(index)));
}

char TxListAccountsModel687eda_RemoveColumnsDefault(void* ptr, int column, int count, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::removeColumns(column, count, *static_cast<QModelIndex*>(parent));
}

char TxListAccountsModel687eda_RemoveRowsDefault(void* ptr, int row, int count, void* parent)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::removeRows(row, count, *static_cast<QModelIndex*>(parent));
}

void TxListAccountsModel687eda_ResetInternalDataDefault(void* ptr)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::resetInternalData();
}

void TxListAccountsModel687eda_RevertDefault(void* ptr)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::revert();
}

struct Moc_PackedList TxListAccountsModel687eda_RoleNamesDefault(void* ptr)
{
	return ({ QHash<int, QByteArray>* tmpValue5a9af8 = new QHash<int, QByteArray>(static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::roleNames()); Moc_PackedList { tmpValue5a9af8, tmpValue5a9af8->size() }; });
}

int TxListAccountsModel687eda_RowCountDefault(void* ptr, void* parent)
{
	Q_UNUSED(ptr);
	Q_UNUSED(parent);

}

char TxListAccountsModel687eda_SetDataDefault(void* ptr, void* index, void* value, int role)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::setData(*static_cast<QModelIndex*>(index), *static_cast<QVariant*>(value), role);
}

char TxListAccountsModel687eda_SetHeaderDataDefault(void* ptr, int section, long long orientation, void* value, int role)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::setHeaderData(section, static_cast<Qt::Orientation>(orientation), *static_cast<QVariant*>(value), role);
}

char TxListAccountsModel687eda_SetItemDataDefault(void* ptr, void* index, void* roles)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::setItemData(*static_cast<QModelIndex*>(index), *static_cast<QMap<int, QVariant>*>(roles));
}

void TxListAccountsModel687eda_SortDefault(void* ptr, int column, long long order)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::sort(column, static_cast<Qt::SortOrder>(order));
}

void* TxListAccountsModel687eda_SpanDefault(void* ptr, void* index)
{
	return ({ QSize tmpValue = static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::span(*static_cast<QModelIndex*>(index)); new QSize(tmpValue.width(), tmpValue.height()); });
}

char TxListAccountsModel687eda_SubmitDefault(void* ptr)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::submit();
}

long long TxListAccountsModel687eda_SupportedDragActionsDefault(void* ptr)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::supportedDragActions();
}

long long TxListAccountsModel687eda_SupportedDropActionsDefault(void* ptr)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::supportedDropActions();
}

void TxListAccountsModel687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::childEvent(static_cast<QChildEvent*>(event));
}

void TxListAccountsModel687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void TxListAccountsModel687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::customEvent(static_cast<QEvent*>(event));
}

void TxListAccountsModel687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::deleteLater();
}

void TxListAccountsModel687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char TxListAccountsModel687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::event(static_cast<QEvent*>(e));
}

char TxListAccountsModel687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void TxListAccountsModel687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<TxListAccountsModel687eda*>(ptr)->QAbstractListModel::timerEvent(static_cast<QTimerEvent*>(event));
}

void TxListCtx687eda_ConnectClicked(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(qint32)>(&TxListCtx687eda::clicked), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(qint32)>(&TxListCtx687eda::Signal_Clicked), static_cast<Qt::ConnectionType>(t));
}

void TxListCtx687eda_DisconnectClicked(void* ptr)
{
	QObject::disconnect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(qint32)>(&TxListCtx687eda::clicked), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(qint32)>(&TxListCtx687eda::Signal_Clicked));
}

void TxListCtx687eda_Clicked(void* ptr, int b)
{
	static_cast<TxListCtx687eda*>(ptr)->clicked(b);
}

struct Moc_PackedString TxListCtx687eda_ShortenAddress(void* ptr)
{
	return ({ QByteArray* tf699cf = new QByteArray(static_cast<TxListCtx687eda*>(ptr)->shortenAddress().toUtf8()); Moc_PackedString { const_cast<char*>(tf699cf->prepend("WHITESPACE").constData()+10), tf699cf->size()-10, tf699cf }; });
}

struct Moc_PackedString TxListCtx687eda_ShortenAddressDefault(void* ptr)
{
	return ({ QByteArray* tde87a5 = new QByteArray(static_cast<TxListCtx687eda*>(ptr)->shortenAddressDefault().toUtf8()); Moc_PackedString { const_cast<char*>(tde87a5->prepend("WHITESPACE").constData()+10), tde87a5->size()-10, tde87a5 }; });
}

void TxListCtx687eda_SetShortenAddress(void* ptr, struct Moc_PackedString shortenAddress)
{
	static_cast<TxListCtx687eda*>(ptr)->setShortenAddress(QString::fromUtf8(shortenAddress.data, shortenAddress.len));
}

void TxListCtx687eda_SetShortenAddressDefault(void* ptr, struct Moc_PackedString shortenAddress)
{
	static_cast<TxListCtx687eda*>(ptr)->setShortenAddressDefault(QString::fromUtf8(shortenAddress.data, shortenAddress.len));
}

void TxListCtx687eda_ConnectShortenAddressChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::shortenAddressChanged), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::Signal_ShortenAddressChanged), static_cast<Qt::ConnectionType>(t));
}

void TxListCtx687eda_DisconnectShortenAddressChanged(void* ptr)
{
	QObject::disconnect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::shortenAddressChanged), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::Signal_ShortenAddressChanged));
}

void TxListCtx687eda_ShortenAddressChanged(void* ptr, struct Moc_PackedString shortenAddress)
{
	static_cast<TxListCtx687eda*>(ptr)->shortenAddressChanged(QString::fromUtf8(shortenAddress.data, shortenAddress.len));
}

struct Moc_PackedString TxListCtx687eda_SelectedSrc(void* ptr)
{
	return ({ QByteArray* tff7988 = new QByteArray(static_cast<TxListCtx687eda*>(ptr)->selectedSrc().toUtf8()); Moc_PackedString { const_cast<char*>(tff7988->prepend("WHITESPACE").constData()+10), tff7988->size()-10, tff7988 }; });
}

struct Moc_PackedString TxListCtx687eda_SelectedSrcDefault(void* ptr)
{
	return ({ QByteArray* te8fe51 = new QByteArray(static_cast<TxListCtx687eda*>(ptr)->selectedSrcDefault().toUtf8()); Moc_PackedString { const_cast<char*>(te8fe51->prepend("WHITESPACE").constData()+10), te8fe51->size()-10, te8fe51 }; });
}

void TxListCtx687eda_SetSelectedSrc(void* ptr, struct Moc_PackedString selectedSrc)
{
	static_cast<TxListCtx687eda*>(ptr)->setSelectedSrc(QString::fromUtf8(selectedSrc.data, selectedSrc.len));
}

void TxListCtx687eda_SetSelectedSrcDefault(void* ptr, struct Moc_PackedString selectedSrc)
{
	static_cast<TxListCtx687eda*>(ptr)->setSelectedSrcDefault(QString::fromUtf8(selectedSrc.data, selectedSrc.len));
}

void TxListCtx687eda_ConnectSelectedSrcChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::selectedSrcChanged), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::Signal_SelectedSrcChanged), static_cast<Qt::ConnectionType>(t));
}

void TxListCtx687eda_DisconnectSelectedSrcChanged(void* ptr)
{
	QObject::disconnect(static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::selectedSrcChanged), static_cast<TxListCtx687eda*>(ptr), static_cast<void (TxListCtx687eda::*)(QString)>(&TxListCtx687eda::Signal_SelectedSrcChanged));
}

void TxListCtx687eda_SelectedSrcChanged(void* ptr, struct Moc_PackedString selectedSrc)
{
	static_cast<TxListCtx687eda*>(ptr)->selectedSrcChanged(QString::fromUtf8(selectedSrc.data, selectedSrc.len));
}

int TxListCtx687eda_TxListCtx687eda_QRegisterMetaType()
{
	return qRegisterMetaType<TxListCtx687eda*>();
}

int TxListCtx687eda_TxListCtx687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<TxListCtx687eda*>(const_cast<const char*>(typeName));
}

int TxListCtx687eda_TxListCtx687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListCtx687eda>();
#else
	return 0;
#endif
}

int TxListCtx687eda_TxListCtx687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<TxListCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int TxListCtx687eda_TxListCtx687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<TxListCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* TxListCtx687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListCtx687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListCtx687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* TxListCtx687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void TxListCtx687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* TxListCtx687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* TxListCtx687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListCtx687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListCtx687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListCtx687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void TxListCtx687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* TxListCtx687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* TxListCtx687eda_NewTxListCtx(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new TxListCtx687eda(static_cast<QWindow*>(parent));
	} else {
		return new TxListCtx687eda(static_cast<QObject*>(parent));
	}
}

void TxListCtx687eda_DestroyTxListCtx(void* ptr)
{
	static_cast<TxListCtx687eda*>(ptr)->~TxListCtx687eda();
}

void TxListCtx687eda_DestroyTxListCtxDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void TxListCtx687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void TxListCtx687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void TxListCtx687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void TxListCtx687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::deleteLater();
}

void TxListCtx687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char TxListCtx687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<TxListCtx687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char TxListCtx687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<TxListCtx687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void TxListCtx687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<TxListCtx687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void ApproveListingCtx687eda_ConnectBack(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::back), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Back), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectBack(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::back), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Back));
}

void ApproveListingCtx687eda_Back(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->back();
}

void ApproveListingCtx687eda_ConnectApprove(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::approve), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Approve), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectApprove(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::approve), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Approve));
}

void ApproveListingCtx687eda_Approve(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->approve();
}

void ApproveListingCtx687eda_ConnectReject(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::reject), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Reject), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectReject(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::reject), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_Reject));
}

void ApproveListingCtx687eda_Reject(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->reject();
}

void ApproveListingCtx687eda_ConnectOnCheckStateChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(qint32, bool)>(&ApproveListingCtx687eda::onCheckStateChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(qint32, bool)>(&ApproveListingCtx687eda::Signal_OnCheckStateChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectOnCheckStateChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(qint32, bool)>(&ApproveListingCtx687eda::onCheckStateChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(qint32, bool)>(&ApproveListingCtx687eda::Signal_OnCheckStateChanged));
}

void ApproveListingCtx687eda_OnCheckStateChanged(void* ptr, int i, char checked)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->onCheckStateChanged(i, checked != 0);
}

void ApproveListingCtx687eda_ConnectTriggerUpdate(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::triggerUpdate), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_TriggerUpdate), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectTriggerUpdate(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::triggerUpdate), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)()>(&ApproveListingCtx687eda::Signal_TriggerUpdate));
}

void ApproveListingCtx687eda_TriggerUpdate(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->triggerUpdate();
}

struct Moc_PackedString ApproveListingCtx687eda_Remote(void* ptr)
{
	return ({ QByteArray* t9605d3 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->remote().toUtf8()); Moc_PackedString { const_cast<char*>(t9605d3->prepend("WHITESPACE").constData()+10), t9605d3->size()-10, t9605d3 }; });
}

struct Moc_PackedString ApproveListingCtx687eda_RemoteDefault(void* ptr)
{
	return ({ QByteArray* t4b8f13 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->remoteDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t4b8f13->prepend("WHITESPACE").constData()+10), t4b8f13->size()-10, t4b8f13 }; });
}

void ApproveListingCtx687eda_SetRemote(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setRemote(QString::fromUtf8(remote.data, remote.len));
}

void ApproveListingCtx687eda_SetRemoteDefault(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setRemoteDefault(QString::fromUtf8(remote.data, remote.len));
}

void ApproveListingCtx687eda_ConnectRemoteChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::remoteChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_RemoteChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectRemoteChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::remoteChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_RemoteChanged));
}

void ApproveListingCtx687eda_RemoteChanged(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->remoteChanged(QString::fromUtf8(remote.data, remote.len));
}

struct Moc_PackedString ApproveListingCtx687eda_Transport(void* ptr)
{
	return ({ QByteArray* t342796 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->transport().toUtf8()); Moc_PackedString { const_cast<char*>(t342796->prepend("WHITESPACE").constData()+10), t342796->size()-10, t342796 }; });
}

struct Moc_PackedString ApproveListingCtx687eda_TransportDefault(void* ptr)
{
	return ({ QByteArray* t433cbf = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->transportDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t433cbf->prepend("WHITESPACE").constData()+10), t433cbf->size()-10, t433cbf }; });
}

void ApproveListingCtx687eda_SetTransport(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setTransport(QString::fromUtf8(transport.data, transport.len));
}

void ApproveListingCtx687eda_SetTransportDefault(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setTransportDefault(QString::fromUtf8(transport.data, transport.len));
}

void ApproveListingCtx687eda_ConnectTransportChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::transportChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_TransportChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectTransportChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::transportChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_TransportChanged));
}

void ApproveListingCtx687eda_TransportChanged(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->transportChanged(QString::fromUtf8(transport.data, transport.len));
}

struct Moc_PackedString ApproveListingCtx687eda_Endpoint(void* ptr)
{
	return ({ QByteArray* t637b7c = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->endpoint().toUtf8()); Moc_PackedString { const_cast<char*>(t637b7c->prepend("WHITESPACE").constData()+10), t637b7c->size()-10, t637b7c }; });
}

struct Moc_PackedString ApproveListingCtx687eda_EndpointDefault(void* ptr)
{
	return ({ QByteArray* t7e670e = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->endpointDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t7e670e->prepend("WHITESPACE").constData()+10), t7e670e->size()-10, t7e670e }; });
}

void ApproveListingCtx687eda_SetEndpoint(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setEndpoint(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveListingCtx687eda_SetEndpointDefault(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setEndpointDefault(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveListingCtx687eda_ConnectEndpointChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::endpointChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_EndpointChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectEndpointChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::endpointChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_EndpointChanged));
}

void ApproveListingCtx687eda_EndpointChanged(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->endpointChanged(QString::fromUtf8(endpoint.data, endpoint.len));
}

struct Moc_PackedString ApproveListingCtx687eda_From(void* ptr)
{
	return ({ QByteArray* t68097f = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->from().toUtf8()); Moc_PackedString { const_cast<char*>(t68097f->prepend("WHITESPACE").constData()+10), t68097f->size()-10, t68097f }; });
}

struct Moc_PackedString ApproveListingCtx687eda_FromDefault(void* ptr)
{
	return ({ QByteArray* t5c8897 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->fromDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t5c8897->prepend("WHITESPACE").constData()+10), t5c8897->size()-10, t5c8897 }; });
}

void ApproveListingCtx687eda_SetFrom(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setFrom(QString::fromUtf8(from.data, from.len));
}

void ApproveListingCtx687eda_SetFromDefault(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setFromDefault(QString::fromUtf8(from.data, from.len));
}

void ApproveListingCtx687eda_ConnectFromChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::fromChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_FromChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectFromChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::fromChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_FromChanged));
}

void ApproveListingCtx687eda_FromChanged(void* ptr, struct Moc_PackedString from)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->fromChanged(QString::fromUtf8(from.data, from.len));
}

struct Moc_PackedString ApproveListingCtx687eda_Message(void* ptr)
{
	return ({ QByteArray* t87291a = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->message().toUtf8()); Moc_PackedString { const_cast<char*>(t87291a->prepend("WHITESPACE").constData()+10), t87291a->size()-10, t87291a }; });
}

struct Moc_PackedString ApproveListingCtx687eda_MessageDefault(void* ptr)
{
	return ({ QByteArray* t124a6d = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->messageDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t124a6d->prepend("WHITESPACE").constData()+10), t124a6d->size()-10, t124a6d }; });
}

void ApproveListingCtx687eda_SetMessage(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setMessage(QString::fromUtf8(message.data, message.len));
}

void ApproveListingCtx687eda_SetMessageDefault(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setMessageDefault(QString::fromUtf8(message.data, message.len));
}

void ApproveListingCtx687eda_ConnectMessageChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::messageChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_MessageChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectMessageChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::messageChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_MessageChanged));
}

void ApproveListingCtx687eda_MessageChanged(void* ptr, struct Moc_PackedString message)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->messageChanged(QString::fromUtf8(message.data, message.len));
}

struct Moc_PackedString ApproveListingCtx687eda_RawData(void* ptr)
{
	return ({ QByteArray* t409f0f = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->rawData().toUtf8()); Moc_PackedString { const_cast<char*>(t409f0f->prepend("WHITESPACE").constData()+10), t409f0f->size()-10, t409f0f }; });
}

struct Moc_PackedString ApproveListingCtx687eda_RawDataDefault(void* ptr)
{
	return ({ QByteArray* t518430 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->rawDataDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t518430->prepend("WHITESPACE").constData()+10), t518430->size()-10, t518430 }; });
}

void ApproveListingCtx687eda_SetRawData(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setRawData(QString::fromUtf8(rawData.data, rawData.len));
}

void ApproveListingCtx687eda_SetRawDataDefault(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setRawDataDefault(QString::fromUtf8(rawData.data, rawData.len));
}

void ApproveListingCtx687eda_ConnectRawDataChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::rawDataChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_RawDataChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectRawDataChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::rawDataChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_RawDataChanged));
}

void ApproveListingCtx687eda_RawDataChanged(void* ptr, struct Moc_PackedString rawData)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->rawDataChanged(QString::fromUtf8(rawData.data, rawData.len));
}

struct Moc_PackedString ApproveListingCtx687eda_Hash(void* ptr)
{
	return ({ QByteArray* t74e19d = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->hash().toUtf8()); Moc_PackedString { const_cast<char*>(t74e19d->prepend("WHITESPACE").constData()+10), t74e19d->size()-10, t74e19d }; });
}

struct Moc_PackedString ApproveListingCtx687eda_HashDefault(void* ptr)
{
	return ({ QByteArray* td3c517 = new QByteArray(static_cast<ApproveListingCtx687eda*>(ptr)->hashDefault().toUtf8()); Moc_PackedString { const_cast<char*>(td3c517->prepend("WHITESPACE").constData()+10), td3c517->size()-10, td3c517 }; });
}

void ApproveListingCtx687eda_SetHash(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setHash(QString::fromUtf8(hash.data, hash.len));
}

void ApproveListingCtx687eda_SetHashDefault(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->setHashDefault(QString::fromUtf8(hash.data, hash.len));
}

void ApproveListingCtx687eda_ConnectHashChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::hashChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_HashChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveListingCtx687eda_DisconnectHashChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::hashChanged), static_cast<ApproveListingCtx687eda*>(ptr), static_cast<void (ApproveListingCtx687eda::*)(QString)>(&ApproveListingCtx687eda::Signal_HashChanged));
}

void ApproveListingCtx687eda_HashChanged(void* ptr, struct Moc_PackedString hash)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->hashChanged(QString::fromUtf8(hash.data, hash.len));
}

int ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType()
{
	return qRegisterMetaType<ApproveListingCtx687eda*>();
}

int ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<ApproveListingCtx687eda*>(const_cast<const char*>(typeName));
}

int ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveListingCtx687eda>();
#else
	return 0;
#endif
}

int ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveListingCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<ApproveListingCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* ApproveListingCtx687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveListingCtx687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveListingCtx687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* ApproveListingCtx687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void ApproveListingCtx687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* ApproveListingCtx687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* ApproveListingCtx687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveListingCtx687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveListingCtx687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveListingCtx687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveListingCtx687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveListingCtx687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveListingCtx687eda_NewApproveListingCtx(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveListingCtx687eda(static_cast<QWindow*>(parent));
	} else {
		return new ApproveListingCtx687eda(static_cast<QObject*>(parent));
	}
}

void ApproveListingCtx687eda_DestroyApproveListingCtx(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->~ApproveListingCtx687eda();
}

void ApproveListingCtx687eda_DestroyApproveListingCtxDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void ApproveListingCtx687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void ApproveListingCtx687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void ApproveListingCtx687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void ApproveListingCtx687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::deleteLater();
}

void ApproveListingCtx687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char ApproveListingCtx687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<ApproveListingCtx687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char ApproveListingCtx687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<ApproveListingCtx687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void ApproveListingCtx687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<ApproveListingCtx687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

void ApproveNewAccountCtx687eda_ConnectClicked(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(qint32)>(&ApproveNewAccountCtx687eda::clicked), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(qint32)>(&ApproveNewAccountCtx687eda::Signal_Clicked), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectClicked(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(qint32)>(&ApproveNewAccountCtx687eda::clicked), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(qint32)>(&ApproveNewAccountCtx687eda::Signal_Clicked));
}

void ApproveNewAccountCtx687eda_Clicked(void* ptr, int b)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->clicked(b);
}

void ApproveNewAccountCtx687eda_ConnectBack(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)()>(&ApproveNewAccountCtx687eda::back), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)()>(&ApproveNewAccountCtx687eda::Signal_Back), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectBack(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)()>(&ApproveNewAccountCtx687eda::back), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)()>(&ApproveNewAccountCtx687eda::Signal_Back));
}

void ApproveNewAccountCtx687eda_Back(void* ptr)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->back();
}

void ApproveNewAccountCtx687eda_ConnectPasswordEdited(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::passwordEdited), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_PasswordEdited), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectPasswordEdited(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::passwordEdited), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_PasswordEdited));
}

void ApproveNewAccountCtx687eda_PasswordEdited(void* ptr, struct Moc_PackedString b)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->passwordEdited(QString::fromUtf8(b.data, b.len));
}

void ApproveNewAccountCtx687eda_ConnectConfirmPasswordEdited(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::confirmPasswordEdited), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_ConfirmPasswordEdited), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectConfirmPasswordEdited(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::confirmPasswordEdited), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_ConfirmPasswordEdited));
}

void ApproveNewAccountCtx687eda_ConfirmPasswordEdited(void* ptr, struct Moc_PackedString b)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->confirmPasswordEdited(QString::fromUtf8(b.data, b.len));
}

struct Moc_PackedString ApproveNewAccountCtx687eda_Remote(void* ptr)
{
	return ({ QByteArray* t3094cf = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->remote().toUtf8()); Moc_PackedString { const_cast<char*>(t3094cf->prepend("WHITESPACE").constData()+10), t3094cf->size()-10, t3094cf }; });
}

struct Moc_PackedString ApproveNewAccountCtx687eda_RemoteDefault(void* ptr)
{
	return ({ QByteArray* t63aa2c = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->remoteDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t63aa2c->prepend("WHITESPACE").constData()+10), t63aa2c->size()-10, t63aa2c }; });
}

void ApproveNewAccountCtx687eda_SetRemote(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setRemote(QString::fromUtf8(remote.data, remote.len));
}

void ApproveNewAccountCtx687eda_SetRemoteDefault(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setRemoteDefault(QString::fromUtf8(remote.data, remote.len));
}

void ApproveNewAccountCtx687eda_ConnectRemoteChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::remoteChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_RemoteChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectRemoteChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::remoteChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_RemoteChanged));
}

void ApproveNewAccountCtx687eda_RemoteChanged(void* ptr, struct Moc_PackedString remote)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->remoteChanged(QString::fromUtf8(remote.data, remote.len));
}

struct Moc_PackedString ApproveNewAccountCtx687eda_Transport(void* ptr)
{
	return ({ QByteArray* t98cdda = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->transport().toUtf8()); Moc_PackedString { const_cast<char*>(t98cdda->prepend("WHITESPACE").constData()+10), t98cdda->size()-10, t98cdda }; });
}

struct Moc_PackedString ApproveNewAccountCtx687eda_TransportDefault(void* ptr)
{
	return ({ QByteArray* t73eeff = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->transportDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t73eeff->prepend("WHITESPACE").constData()+10), t73eeff->size()-10, t73eeff }; });
}

void ApproveNewAccountCtx687eda_SetTransport(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setTransport(QString::fromUtf8(transport.data, transport.len));
}

void ApproveNewAccountCtx687eda_SetTransportDefault(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setTransportDefault(QString::fromUtf8(transport.data, transport.len));
}

void ApproveNewAccountCtx687eda_ConnectTransportChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::transportChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_TransportChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectTransportChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::transportChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_TransportChanged));
}

void ApproveNewAccountCtx687eda_TransportChanged(void* ptr, struct Moc_PackedString transport)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->transportChanged(QString::fromUtf8(transport.data, transport.len));
}

struct Moc_PackedString ApproveNewAccountCtx687eda_Endpoint(void* ptr)
{
	return ({ QByteArray* t7d2dd9 = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->endpoint().toUtf8()); Moc_PackedString { const_cast<char*>(t7d2dd9->prepend("WHITESPACE").constData()+10), t7d2dd9->size()-10, t7d2dd9 }; });
}

struct Moc_PackedString ApproveNewAccountCtx687eda_EndpointDefault(void* ptr)
{
	return ({ QByteArray* t43ee71 = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->endpointDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t43ee71->prepend("WHITESPACE").constData()+10), t43ee71->size()-10, t43ee71 }; });
}

void ApproveNewAccountCtx687eda_SetEndpoint(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setEndpoint(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveNewAccountCtx687eda_SetEndpointDefault(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setEndpointDefault(QString::fromUtf8(endpoint.data, endpoint.len));
}

void ApproveNewAccountCtx687eda_ConnectEndpointChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::endpointChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_EndpointChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectEndpointChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::endpointChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_EndpointChanged));
}

void ApproveNewAccountCtx687eda_EndpointChanged(void* ptr, struct Moc_PackedString endpoint)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->endpointChanged(QString::fromUtf8(endpoint.data, endpoint.len));
}

struct Moc_PackedString ApproveNewAccountCtx687eda_Password(void* ptr)
{
	return ({ QByteArray* tcbaaee = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->password().toUtf8()); Moc_PackedString { const_cast<char*>(tcbaaee->prepend("WHITESPACE").constData()+10), tcbaaee->size()-10, tcbaaee }; });
}

struct Moc_PackedString ApproveNewAccountCtx687eda_PasswordDefault(void* ptr)
{
	return ({ QByteArray* t2050de = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->passwordDefault().toUtf8()); Moc_PackedString { const_cast<char*>(t2050de->prepend("WHITESPACE").constData()+10), t2050de->size()-10, t2050de }; });
}

void ApproveNewAccountCtx687eda_SetPassword(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setPassword(QString::fromUtf8(password.data, password.len));
}

void ApproveNewAccountCtx687eda_SetPasswordDefault(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setPasswordDefault(QString::fromUtf8(password.data, password.len));
}

void ApproveNewAccountCtx687eda_ConnectPasswordChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::passwordChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_PasswordChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectPasswordChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::passwordChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_PasswordChanged));
}

void ApproveNewAccountCtx687eda_PasswordChanged(void* ptr, struct Moc_PackedString password)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->passwordChanged(QString::fromUtf8(password.data, password.len));
}

struct Moc_PackedString ApproveNewAccountCtx687eda_ConfirmPassword(void* ptr)
{
	return ({ QByteArray* t47d33e = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->confirmPassword().toUtf8()); Moc_PackedString { const_cast<char*>(t47d33e->prepend("WHITESPACE").constData()+10), t47d33e->size()-10, t47d33e }; });
}

struct Moc_PackedString ApproveNewAccountCtx687eda_ConfirmPasswordDefault(void* ptr)
{
	return ({ QByteArray* te8e3d4 = new QByteArray(static_cast<ApproveNewAccountCtx687eda*>(ptr)->confirmPasswordDefault().toUtf8()); Moc_PackedString { const_cast<char*>(te8e3d4->prepend("WHITESPACE").constData()+10), te8e3d4->size()-10, te8e3d4 }; });
}

void ApproveNewAccountCtx687eda_SetConfirmPassword(void* ptr, struct Moc_PackedString confirmPassword)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setConfirmPassword(QString::fromUtf8(confirmPassword.data, confirmPassword.len));
}

void ApproveNewAccountCtx687eda_SetConfirmPasswordDefault(void* ptr, struct Moc_PackedString confirmPassword)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->setConfirmPasswordDefault(QString::fromUtf8(confirmPassword.data, confirmPassword.len));
}

void ApproveNewAccountCtx687eda_ConnectConfirmPasswordChanged(void* ptr, long long t)
{
	QObject::connect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::confirmPasswordChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_ConfirmPasswordChanged), static_cast<Qt::ConnectionType>(t));
}

void ApproveNewAccountCtx687eda_DisconnectConfirmPasswordChanged(void* ptr)
{
	QObject::disconnect(static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::confirmPasswordChanged), static_cast<ApproveNewAccountCtx687eda*>(ptr), static_cast<void (ApproveNewAccountCtx687eda::*)(QString)>(&ApproveNewAccountCtx687eda::Signal_ConfirmPasswordChanged));
}

void ApproveNewAccountCtx687eda_ConfirmPasswordChanged(void* ptr, struct Moc_PackedString confirmPassword)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->confirmPasswordChanged(QString::fromUtf8(confirmPassword.data, confirmPassword.len));
}

int ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType()
{
	return qRegisterMetaType<ApproveNewAccountCtx687eda*>();
}

int ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<ApproveNewAccountCtx687eda*>(const_cast<const char*>(typeName));
}

int ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveNewAccountCtx687eda>();
#else
	return 0;
#endif
}

int ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ApproveNewAccountCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

int ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterUncreatableType(char* uri, int versionMajor, int versionMinor, char* qmlName, struct Moc_PackedString message)
{
#ifdef QT_QML_LIB
	return qmlRegisterUncreatableType<ApproveNewAccountCtx687eda>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName), QString::fromUtf8(message.data, message.len));
#else
	return 0;
#endif
}

void* ApproveNewAccountCtx687eda___children_atList(void* ptr, int i)
{
	return ({QObject * tmp = static_cast<QList<QObject *>*>(ptr)->at(i); if (i == static_cast<QList<QObject *>*>(ptr)->size()-1) { static_cast<QList<QObject *>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveNewAccountCtx687eda___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveNewAccountCtx687eda___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>();
}

void* ApproveNewAccountCtx687eda___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(({QByteArray tmp = static_cast<QList<QByteArray>*>(ptr)->at(i); if (i == static_cast<QList<QByteArray>*>(ptr)->size()-1) { static_cast<QList<QByteArray>*>(ptr)->~QList(); free(ptr); }; tmp; }));
}

void ApproveNewAccountCtx687eda___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* ApproveNewAccountCtx687eda___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>();
}

void* ApproveNewAccountCtx687eda___findChildren_atList(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveNewAccountCtx687eda___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveNewAccountCtx687eda___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveNewAccountCtx687eda___findChildren_atList3(void* ptr, int i)
{
	return ({QObject* tmp = static_cast<QList<QObject*>*>(ptr)->at(i); if (i == static_cast<QList<QObject*>*>(ptr)->size()-1) { static_cast<QList<QObject*>*>(ptr)->~QList(); free(ptr); }; tmp; });
}

void ApproveNewAccountCtx687eda___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ApproveNewAccountCtx687eda___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>();
}

void* ApproveNewAccountCtx687eda_NewApproveNewAccountCtx(void* parent)
{
	if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new ApproveNewAccountCtx687eda(static_cast<QWindow*>(parent));
	} else {
		return new ApproveNewAccountCtx687eda(static_cast<QObject*>(parent));
	}
}

void ApproveNewAccountCtx687eda_DestroyApproveNewAccountCtx(void* ptr)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->~ApproveNewAccountCtx687eda();
}

void ApproveNewAccountCtx687eda_DestroyApproveNewAccountCtxDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

void ApproveNewAccountCtx687eda_ChildEventDefault(void* ptr, void* event)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void ApproveNewAccountCtx687eda_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void ApproveNewAccountCtx687eda_CustomEventDefault(void* ptr, void* event)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void ApproveNewAccountCtx687eda_DeleteLaterDefault(void* ptr)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::deleteLater();
}

void ApproveNewAccountCtx687eda_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

char ApproveNewAccountCtx687eda_EventDefault(void* ptr, void* e)
{
	return static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char ApproveNewAccountCtx687eda_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}



void ApproveNewAccountCtx687eda_TimerEventDefault(void* ptr, void* event)
{
	static_cast<ApproveNewAccountCtx687eda*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}

#include "moc_moc.h"
