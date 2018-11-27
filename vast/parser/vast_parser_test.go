package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	data := `
	<VAST xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="vast.xsd" version="3.0">
	blahblah
	<Ad id="31423" sequence="0">
	<InLine>
	<AdSystem version="1.0">1.0</AdSystem>
	<AdTitle/>
	<Description/>
	<Error>
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=err&tpy=m&eid=5bedb41b8df6bf0001c5abd2
	]]>
	</Error>
	<Impression>
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=imp&tpy=m&eid=5bedb41b8df6bf0001c5abd3
	]]>
	</Impression>
	<Creatives>
	<Creative id="70914">
	<Linear>
	<Duration>00:00:53</Duration>
	<TrackingEvents>
	<Tracking event="creativeView">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=cv&eid=5bedb41b8df6bf0001c5abd6
	]]>
	</Tracking>
	<Tracking event="start">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=st&eid=5bedb41b8df6bf0001c5abd7
	]]>
	</Tracking>
	<Tracking event="midpoint">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=mp&eid=5bedb41b8df6bf0001c5abd8
	]]>
	</Tracking>
	<Tracking event="firstQuartile">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=fq&eid=5bedb41b8df6bf0001c5abd9
	]]>
	</Tracking>
	<Tracking event="thirdQuartile">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=tq&eid=5bedb41b8df6bf0001c5abda
	]]>
	</Tracking>
	<Tracking event="complete">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=cp&eid=5bedb41b8df6bf0001c5abdb
	]]>
	</Tracking>
	<Tracking event="mute">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=mt&eid=5bedb41b8df6bf0001c5abdc
	]]>
	</Tracking>
	<Tracking event="unmute">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=um&eid=5bedb41b8df6bf0001c5abdd
	]]>
	</Tracking>
	<Tracking event="pause">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=ps&eid=5bedb41b8df6bf0001c5abde
	]]>
	</Tracking>
	<Tracking event="rewind">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=rw&eid=5bedb41b8df6bf0001c5abdf
	]]>
	</Tracking>
	<Tracking event="resume">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=rs&eid=5bedb41b8df6bf0001c5abe0
	]]>
	</Tracking>
	<Tracking event="fullscreen">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=fs&eid=5bedb41b8df6bf0001c5abe1
	]]>
	</Tracking>
	<Tracking event="expand">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=ep&eid=5bedb41b8df6bf0001c5abe2
	]]>
	</Tracking>
	<Tracking event="collapse">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=cl&eid=5bedb41b8df6bf0001c5abe3
	]]>
	</Tracking>
	<Tracking event="acceptInvitation">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=ai&eid=5bedb41b8df6bf0001c5abe4
	]]>
	</Tracking>
	<Tracking event="close">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=tge&tpy=cs&eid=5bedb41b8df6bf0001c5abe5
	]]>
	</Tracking>
	</TrackingEvents>
	<VideoClicks>
	<ClickThrough>
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=cli&tpy=th&eid=5bedb41b8df6bf0001c5abd5
	]]>
	</ClickThrough>
	</VideoClicks>
	<MediaFiles>
	<MediaFile id="5bd335f6b13d2e15dae28257" delivery="progressive" type="video/mp4" maintainAspectRatio="false" scalable="true" bitrate="4655630" width="1920" height="1080" apiFramework="">
	<![CDATA[
	http://vipk-cache.cdnvideo.ru/vi/test_td/5bd335f6b13d2e15dae28257.mp4
	]]>
	</MediaFile>
	</MediaFiles>
	</Linear>
	</Creative>
	</Creatives>
	<Extensions>
	<Extension type="CustomTracking">
	<Tracking event="onVastLoad">
	<![CDATA[
	https://ads.nsc-lab.io/ads/vast/event?adid=31423&knd=ext&tpy=m&eid=5bedb41b8df6bf0001c5abd4
	]]>
	</Tracking>
	</Extension>
	<Extension type="isClickable">
	<![CDATA[ 1 ]]>
	</Extension>
	</Extensions>
	</InLine>
	</Ad>
	</VAST>
	`

	v := Parse([]byte(data))
	t.Logf("Impression: %v\n", strings.TrimSpace(v.Impression))
	for _, tracking := range v.Creative[0].TrackingEvents {
		t.Logf("TrackingEvents: %v\n", strings.TrimSpace(tracking.URL))
	}

	t.Logf("Creative ID: %v\n", v.Creative[0].ID)
}
