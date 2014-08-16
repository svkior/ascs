import QtQuick 2.1


Rectangle {
	id: screen
	width: 800; height: 600
	color: "black"
	Text {
		id: vvod
		x: 8
		y: 10
		color: "white"
		font.pointSize: 20
		text: "Ввод:"
	}
	

	TextInput {
		id: inp1
		objectName: "inputtext1"
		width: 240
		x : 120
		y : 10
		text: "www.ya.ru:80"
		font.pointSize: 20
		color: "steelblue"
		focus: true
		onAccepted: bridge.handleClick()
		Component.onCompleted: inp1.selectAll()
	}
	

	ListView{
		id: listView1
		objectName: "logview"
		anchors.right: parent.right
		anchors.left: parent.left
		anchors.leftMargin: 10
		anchors.bottom: bottomBar.top
		anchors.bottomMargin: 30
		anchors.top: vvod.bottom
		anchors.topMargin: 30
		//anchors.top: parent.top 
		currentIndex: -1
		onCountChanged: {
			listView1.positionViewAtIndex(listView1.count -1, listView1.End)
		}
		model: logs.len
		delegate: Row {
			id: row1
			height: 30
			Rectangle { 
				id: boxrect
				width: 40; height: 20; color: logs.color(index); radius: 20
				Text {
					anchors.centerIn: parent
					text: index
				}
			}
			Rectangle { width: 20; height:40; color: "black"}
			Text {
				text: logs.record(index)
				color: logs.color(index)
				anchors.verticalCenter: boxrect.verticalCenter
				font.bold: true
			}
		}
	}


	Text {
		id: bottomBar
		anchors.horizontalCenter: parent.horizontalCenter
		anchors.bottom: parent.bottom
		text: "UDP DayTime Клиент и Сервер"
		color: "steelblue"
		font.pointSize: 20
	}

}