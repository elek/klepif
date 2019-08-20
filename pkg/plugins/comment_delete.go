package plugins

//
//type CommentFix struct {
//	DefaultPlugin
//}
//
//func (*CommentFix) HandleComment(cc *ClientContext, pr *github.PullRequest, comments []*github.IssueComment) error {
//	fmt.Println(pr.GetTitle())
//	fmt.Println(pr.GetNumber())
//	for _, comment := range comments {
//		if comment.GetBody() == "Can one of the admins verify this patch?" && comment.User.GetLogin() == "elek" {
//			fmt.Println("this is bad")
//			fmt.Println(comment.GetBody())
//			_, err := cc.Client.Issues.DeleteComment(cc.Ctx, "apache", "hadoop", comment.GetID())
//			if err != nil {
//				return err
//			}
//
//		}
//	}
//	fmt.Println()
//	return nil
//}
